// Copyright (c) 2021 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The tailscale-synology tool generates Tailscale Synology SPK packages.
package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing/fstest"
	"time"

	"github.com/ulikunitz/xz"
)

// srcFS are all the misc files that need to be in the SPK package.
// We embed them into the binary so the tailscale-synology tool can
// be run anywhere. Plus they're small.
//
//go:embed src/*
var srcFS embed.FS

var (
	dsmVersion    = flag.String("dsm-version", "7", `DSM version(s) to build: 6, 7, or "all"`)
	goarch        = flag.String("goarch", "amd64", `GOARCH to build package(s) for, "all`)
	compress      = flag.String("compress", "speed", `compression option: "speed" or "size"; empty means automatic where local builds are fast but big`)
	packageCenter = flag.Bool("for-package-center", false, `build for the package center`)

	srcDir = flag.String("source", ".", "path to tailscale.com's go.mod directory root")
	output = flag.String("o", "", "output directory or path; if a directory, files are written there. If it ends in *.spk, the spk is written to that name")
)

// synPlat maps from GOARCH (or GOARCH/GOARM) to the Synology platform name(s).
//
// architecture taken from:
// https://github.com/SynoCommunity/spksrc/wiki/Synology-and-SynoCommunity-Package-Architectures
// https://github.com/SynologyOpenSource/pkgscripts-ng/tree/master/include platform.<PLATFORM> files
var synPlat = map[string][]string{
	"amd64": {"x86_64"},
	"386":   {"i686"},
	"arm64": {"armv8"},
	"arm/5": {"armv5", "88f6281", "88f6282"},
	"arm/7": {"armv7", "alpine", "armada370", "armada375", "armada38x", "armadaxp", "comcerto2k", "monaco", "hi3535"},
}

func main() {
	flag.Parse()

	switch *compress {
	case "size", "speed":
	default:
		log.Fatalf("invalid --compress value %q", *compress)
	}

	var dsms []int
	switch *dsmVersion {
	default:
		log.Fatalf("invalid --dsm-version %q", *dsmVersion)
	case "6", "7":
		dsms = append(dsms, int((*dsmVersion)[0]-'0'))
	case "all":
		dsms = append(dsms, 6, 7)
	}
	var goarchs []string
	if *goarch == "all" {
		for goarch := range synPlat {
			goarchs = append(goarchs, goarch)
		}
	} else {
		if vv := synPlat[*goarch]; len(vv) == 0 {
			log.Fatalf("unknown --goarch value %q", *goarch)
		}
		goarchs = append(goarchs, *goarch)
	}
	if len(goarchs) > 1 || len(dsms) > 1 {
		fi, err := os.Stat(*output)
		if err != nil {
			log.Fatal(err)
		}
		if !fi.IsDir() {
			log.Fatalf("%q is not a dir", *output)
		}
	}

	dv, err := getDistVars(*srcDir)
	if err != nil {
		log.Fatal(err)
	}

	commitTime, err := readCommitTime(*srcDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, goarch := range goarchs {
		param := spkParams{
			createTime:       commitTime,
			version:          dv.MajorMinorPatch,
			spkBuildBase:     dv.SPKBuild,
			goarch:           goarch,
			forPackageCenter: *packageCenter,
			srcDir:           *srcDir,
		}

		for _, dsm := range dsms {
			if err := genArchSPKs(param, dsm, synPlat[goarch]); err != nil {
				log.Fatal(err)
			}
		}
	}
}

type spkParams struct {
	createTime       time.Time
	spkBuildBase     int    // derived from the short version.
	version          string // "1.18.1", Tailscale short version
	goarch           string // "amd64"
	forPackageCenter bool

	// srcDir, if non-empty, means to use the "go" command to build
	// the tailscale and tailscaled binaries in the srcDir directory
	// instead of downloading them from pkgs.tailscale.com.
	srcDir string
}

func readCommitTime(dir string) (time.Time, error) {
	cmd := exec.Command("git", "show", "--format=%ct", "--quiet")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return time.Time{}, err
	}
	unixString := strings.TrimSpace(string(out))
	unix, err := strconv.ParseInt(unixString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(unix, 0), nil
}

func (p spkParams) spkBuild(dsm int) int {
	return 10*p.spkBuildBase + dsm
}

func (p spkParams) versionDashBuild(dsm int) string {
	return fmt.Sprintf("%v-%v", p.version, p.spkBuild(dsm))
}

// filename returns the SPK's base filename.
func (p spkParams) filename(dsm int, synoArch string) string {
	return fmt.Sprintf("tailscale-%v-%v-dsm%v.spk",
		synoArch,
		p.versionDashBuild(dsm),
		dsm)
}

// outputFile returns the path to write the SPK out to based
// on the -o flag value.
func (p spkParams) outputFile(dsm int, synoArch string) string {
	base := p.filename(dsm, synoArch)
	if *output == "" {
		return base
	}
	if strings.HasSuffix(*output, ".spk") {
		return *output
	}
	return filepath.Join(*output, base)
}

func (p spkParams) goEnv() []string {
	const armPrefix = "arm/"
	goarch := p.goarch
	var goarm string
	if strings.HasPrefix(goarch, armPrefix) {
		goarm = strings.TrimPrefix(goarch, armPrefix)
		goarch = "arm"
	}
	env := append(os.Environ(),
		"CGO_ENABLED=0",
		"GOOS=linux",
		"GOARCH="+goarch,
		"GOARM="+goarm,
	)
	return env
}

func genSPK(param spkParams, dsm int, synoArch string, extractedSize int64, privFile string, innerPackage []byte) error {
	out := param.outputFile(dsm, synoArch)
	log.Printf("Generating %v ...", filepath.Base(out))
	w, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer w.Close()
	return writeTar(w,
		file("INFO", getInfo(param, dsm, synoArch, extractedSize)),
		file("PACKAGE_ICON.PNG", static("PACKAGE_ICON.PNG")),
		file("PACKAGE_ICON_256.PNG", static("PACKAGE_ICON_256.PNG")),
		file("Tailscale.sc", static("Tailscale.sc")),
		dir("conf/", param.createTime),
		file("conf/resource", static("resource")),
		file("conf/privilege", static(privFile)),
		file("package.tgz", memFile(innerPackage, 0644, param.createTime)),
		dir("scripts/", param.createTime),
		file("scripts/start-stop-status", static("scripts/start-stop-status")),
		file("scripts/postupgrade", static("scripts/postupgrade")),
		file("scripts/preupgrade", static("scripts/preupgrade")),
	)
}

// genArchSPKs generates SPKs for a particular DSM version.
// This function generates one spk per synoArch.
// SPKs are nested tarballs.
// The outer tar file is uncompressed and contains the minimal
// metadata. The main contents are in the outer tar's "package.tgz"
// entry, which is gzip or xz compressed.
func genArchSPKs(param spkParams, dsm int, synoArchs []string) error {
	var innerPkgTgz bytes.Buffer
	extractedSize, err := genInnerPackageTgz(&innerPkgTgz, param, dsm)
	if err != nil {
		return err
	}

	privFile := fmt.Sprintf("privilege-dsm%d", dsm)
	if param.forPackageCenter {
		privFile += ".priv"
	}

	innerPackageBytes := innerPkgTgz.Bytes()

	for _, synoArch := range synoArchs {
		if err := genSPK(param, dsm, synoArch, extractedSize, privFile, innerPackageBytes); err != nil {
			return err
		}
	}
	return nil
}

func genInnerPackageTgz(w io.Writer, param spkParams, dsm int) (extractedSize int64, err error) {
	var wc io.WriteCloser
	switch *compress {
	case "speed":
		wc, err = gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			return 0, err
		}
	case "size":
		wc, err = xz.NewWriter(w)
		if err != nil {
			return 0, err
		}
	}
	err = writeTar(io.MultiWriter(writeByteCounter{&extractedSize}, wc),
		dir("bin/", param.createTime),
		file("bin/tailscaled", buildBin("tailscaled", param)),
		file("bin/tailscale", buildBin("tailscale", param)),
		dir("conf/", param.createTime),
		file("conf/Tailscale.sc", static("Tailscale.sc")),
		file("conf/logrotate.conf", static("logrotate-dsm"+strconv.Itoa(dsm))),
		dir("ui/", param.createTime),
		file("ui/PACKAGE_ICON_256.PNG", static("PACKAGE_ICON_256.PNG")),
		file("ui/config", static("config")), // TODO: this has "1.8.3" hard-coded in it; why? what is it? bug?
		file("ui/index.cgi", static("index.cgi")),
	)
	if err != nil {
		return 0, err
	}
	return extractedSize, wc.Close()
}

type writeByteCounter struct{ sum *int64 }

func (w writeByteCounter) Write(p []byte) (int, error) {
	*w.sum += int64(len(p))
	return len(p), nil
}

type tarEntry func(*tar.Writer) error

func writeTar(w io.Writer, ents ...tarEntry) error {
	tw := tar.NewWriter(w)
	for _, ent := range ents {
		if err := ent(tw); err != nil {
			return err
		}
	}
	return tw.Close()
}

func dir(name string, modTime time.Time) tarEntry {
	if !strings.HasSuffix(name, "/") {
		name += "/"
	}
	return func(tw *tar.Writer) error {
		return tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeDir,
			Name:     name,
			Mode:     0755,
			ModTime:  modTime.Truncate(time.Second),
			Uname:    "tailscale",
			Gname:    "tailscale",
		})
	}
}

type fileOpener func() (fs.File, error)

func file(name string, open fileOpener) tarEntry {
	return func(tw *tar.Writer) error {
		f, err := open()
		if err != nil {
			return err
		}
		defer f.Close()
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(&tar.Header{
			Typeflag: tar.TypeReg,
			Name:     name,
			Mode:     int64(fi.Mode().Perm()),
			ModTime:  fi.ModTime().Truncate(time.Second),
			Size:     fi.Size(),
		}); err != nil {
			return err
		}
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		return nil
	}
}

func static(name string) fileOpener {
	return func() (fs.File, error) {
		return srcFS.Open("src/" + name)
	}
}

func memFile(data []byte, mode fs.FileMode, modTime time.Time) fileOpener {
	// Round the time down to a second so a tar entry can't be
	// in the future for sub-second on filesystems without
	// sub-sec resolution.
	modTime = modTime.Truncate(time.Second)
	return func() (fs.File, error) {
		fs := fstest.MapFS{"foo": &fstest.MapFile{
			Data:    data,
			Mode:    mode,
			ModTime: modTime,
		}}
		return fs.Open("foo")
	}
}

func errLater(err error) fileOpener {
	return func() (fs.File, error) {
		return nil, err
	}
}

// compileGoBinary compiles the binary to a temp file and returns the path.
// Cleaning up the file is the responsibility of caller.
// TODO: This is the same as the one github.com/tailscale/mkctr, share somehow?
func compileGoBinary(dir, gopath string, env []string, ldflags, gotags string) (string, error) {
	f, err := os.CreateTemp("", "out")
	if err != nil {
		return "", err
	}
	out := f.Name()
	if err := f.Close(); err != nil {
		return "", err
	}
	args := []string{
		"build",
		"-v",
		"-trimpath",
	}
	if len(gotags) > 0 {
		args = append(args, "--tags="+gotags)
	}
	if len(ldflags) > 0 {
		args = append(args, "--ldflags="+ldflags)
	}
	args = append(args,
		"-o="+out,
		gopath,
	)
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return out, nil
}

// buildBin builds baseProg ("tailscale" or "tailscaled")
// and returns a fileOpener for it.
func buildBin(baseProg string, param spkParams) fileOpener {
	vars, err := getDistVars(param.srcDir)
	if err != nil {
		return errLater(err)
	}
	ldflags := "-X tailscale.com/version.Long=" + vars.Long + " " +
		"-X tailscale.com/version.Short=" + vars.MajorMinorPatch + " " +
		"-X tailscale.com/version.GitCommit=" + vars.GitHash
	name, err := compileGoBinary(param.srcDir, "tailscale.com/cmd/"+baseProg, param.goEnv(), ldflags, "")
	if err != nil {
		return errLater(err)
	}
	data, err := os.ReadFile(name)
	if err != nil {
		return errLater(err)
	}
	if err := os.Remove(name); err != nil {
		return errLater(err)
	}
	return memFile(data, 0755, param.createTime)
}

type DistVars struct {
	MajorMinor      string // "1.21"
	MajorMinorPatch string // "1.21.17"
	Long            string // "1.21.17-tb4f817065"
	GitHash         string // "b4f8170657cde2a3a21ffee46c9dd028e400fb0f"
	SPKBuild        int    // 210017
}

func getDistVars(dir string) (v DistVars, err error) {
	cmd := exec.Command("./build_dist.sh", "shellvars")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return v, err
	}
	bs := bufio.NewScanner(bytes.NewReader(out))
	for bs.Scan() {
		k, qv, ok := stringsCut(strings.TrimSpace(bs.Text()), "=")
		if !ok {
			continue
		}
		var sp *string
		switch k {
		case "VERSION_MINOR":
			sp = &v.MajorMinor
		case "VERSION_SHORT":
			sp = &v.MajorMinorPatch
		case "VERSION_LONG":
			sp = &v.Long
		case "VERSION_GIT_HASH":
			sp = &v.GitHash
		}
		if sp != nil {
			*sp, err = strconv.Unquote(qv)
			if err != nil {
				return v, err
			}
		}
	}
	if err := bs.Err(); err != nil {
		return v, err
	}
	parts := strings.Split(v.MajorMinorPatch, ".")
	if len(parts) != 3 {
		return v, fmt.Errorf("unexpected version: %v", v.MajorMinorPatch)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return v, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return v, err
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return v, err
	}
	v.SPKBuild = (major-1)*1e6 + minor*1e3 + patch
	return v, nil
}

// getInfo returns a fileOpener for the top-level INFO file.
// See genInfo.
func getInfo(param spkParams, dsm int, synoArch string, extractedSize int64) fileOpener {
	data, err := genInfo(param, dsm, synoArch, extractedSize)
	if err != nil {
		return func() (fs.File, error) { return nil, err }
	}
	return memFile(data, 0644, param.createTime)
}

// genInfo returns the outer tar's INFO file, which looks like:
/*
package="Tailscale"
version="1.16.2-2013"
arch="x86_64"
description="Connect all your devices using WireGuard, without the hassle."
displayname="Tailscale"
maintainer="Tailscale, Inc."
maintainer_url="https://github.com/tailscale/tailscale-synology"
create_time="20211103-21:01:18"
dsmuidir="ui"
dsmappname="SYNO.SDS.Tailscale"
startstop_restart_services="nginx"
os_min_ver="7.0-40000"
os_max_ver=""
extractsize="42368"
*/
func genInfo(param spkParams, dsm int, synoArch string, extractedSize int64) ([]byte, error) {
	var buf bytes.Buffer
	add := func(k, v string) {
		fmt.Fprintf(&buf, "%s=%q\n", k, v)
	}
	add("package", "Tailscale")
	add("version", param.versionDashBuild(dsm))
	add("arch", synoArch)
	add("description", "Connect all your devices using WireGuard, without the hassle.")
	add("displayname", "Tailscale")
	add("maintainer", "Tailscale, Inc.")
	add("maintainer_url", "https://github.com/tailscale/tailscale-synology")
	add("create_time", param.createTime.Format("20060102-15:04:05"))
	add("dsmuidir", "ui")
	add("dsmappname", "SYNO.SDS.Tailscale")
	add("startstop_restart_services", "nginx")
	switch dsm {
	case 6:
		add("os_min_ver", "6.0.1-7445")
		add("os_max_ver", "7.0-40000")
	case 7:
		add("os_min_ver", "7.0-40000")
		add("os_max_ver", "")
	default:
		return nil, fmt.Errorf("unsupported DSM version '%v'", dsm)
	}
	add("extractsize", fmt.Sprintf("%v", extractedSize>>10)) // in KiB
	return buf.Bytes(), nil

}

// stringsCut is strings.Cut from Go 1.18.
func stringsCut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
