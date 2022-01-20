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
	"errors"
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
	version    = flag.String("version", "", `version number to build; "build" means to build locally; see --source`)
	spkBuild   = flag.Int("spk-build", 15, `SPK build number; needs to be monotonically increasing regardless of the --version`)
	dsmVersion = flag.String("dsm-version", "7", `DSM version(s) to build: 6, 7, or "all"`)
	goarch     = flag.String("goarch", "amd64", `GOARCH to build package(s) for, "all`)
	compress   = flag.String("compress", "", `compression option: "speed" or "size"; empty means automatic where local builds are fast but big`)

	srcDir = flag.String("source", ".", "path to tailscale.com's go.mod directory root, when using --version=build")
	output = flag.String("o", "", "output directory or path; if a directory, files are written there. If it ends in *.spk, the final is written to that name in --version=build mode")
)

// synPlat maps from GOARCH (or GOARCH/GOARM) to the Synology platform name(s).
//
// architecture taken from:
// https://github.com/SynoCommunity/spksrc/wiki/Synology-and-SynoCommunity-Package-Architectures
// https://github.com/SynologyOpenSource/pkgscripts-ng/tree/master/include platform.<PLATFORM> files
var synPlat = map[string][]string{
	"amd64": []string{"x86_64"},
	"386":   []string{"i686"},
	"arm64": []string{"armv8"},
	"arm/5": []string{"armv5", "88f6281", "88f6282"},
	"arm/7": []string{"armv7", "alpine", "armada370", "armada375", "armada38x", "armadaxp", "comcerto2k", "monaco", "hi3535"},
}

func main() {
	flag.Parse()

	var doBuild bool
	switch *version {
	case "":
		log.Fatalf("no --version")
	case "build":
		doBuild = true
	default:
		log.Fatalf("TODO: only --version=build is currently supported")
	}
	switch *compress {
	case "size", "speed":
	case "":
		*compress = "size"
		if doBuild {
			*compress = "speed"
		}
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
		if doBuild {
			log.Fatalf("invalid --dsm=all in --version=build mode")
		}
		dsms = append(dsms, 6, 7)
	}
	if *goarch == "all" && doBuild {
		log.Fatalf("invalid --goarch=all in --version=build mode")
	}
	if strings.HasSuffix(*output, ".spk") && !doBuild {
		log.Fatalf("-o value of *.spk only supported in --version=build mode")
	}

	shortVer := *version
	if doBuild {
		var err error
		shortVer, err = getShortVer(*srcDir)
		if err != nil {
			log.Fatal(err)
		}
	}
	var synoArch string
	if vv := synPlat[*goarch]; len(vv) > 0 {
		synoArch = vv[0]
	} else {
		log.Fatalf("unknown --goarch value %q", *goarch)
	}
	param := spkParams{
		createTime:       time.Now(),
		version:          shortVer,
		spkBuildBase:     *spkBuild,
		goarch:           *goarch,
		synoArch:         synoArch,
		dsm:              dsms[0],
		forPackageCenter: false,
		srcDir:           *srcDir,
	}

	file := param.filename()
	out := param.outputFile()
	log.Printf("Generating %v ...", file)
	var buf bytes.Buffer
	if err := genSPK(&buf, param); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(out, buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
	if out != file {
		log.Printf("Wrote %v as %v", file, out)
	}
}

type spkParams struct {
	createTime       time.Time
	spkBuildBase     int    // montonically increasing; the 2000 is added later
	version          string // "1.18.1", Tailscale short version
	goarch           string // "amd64"
	synoArch         string // "x86_64", etc
	dsm              int    // 6, 7
	forPackageCenter bool

	// srcDir, if non-empty, means to use the "go" command to build
	// the tailscale and tailscaled binaries in the srcDir directory
	// instead of downloading them from pkgs.tailscale.com.
	srcDir string
}

func (p spkParams) spkBuild() int {
	if p.dsm == 7 {
		return 2000 + p.spkBuildBase
	}
	return p.spkBuildBase
}

func (p spkParams) spkVersion() string {
	return fmt.Sprintf("%v-%v", p.version, p.spkBuild())
}

// filename returns the SPK's base filename.
func (p spkParams) filename() string {
	return fmt.Sprintf("tailscale-%v-%v-dsm%v.spk",
		p.synoArch,
		p.spkVersion(),
		p.dsm)
}

// outputFile returns the path to write the SPK out to based
// on the -o flag value.
func (p spkParams) outputFile() string {
	base := p.filename()
	if *output == "" {
		return base
	}
	if strings.HasSuffix(*output, ".spk") {
		return *output
	}
	return filepath.Join(*output, base)
}

func (p spkParams) goEnv() []string {
	return append(os.Environ(),
		"GOOS=linux",
		"GOARCH="+p.goarch,
		// TODO: add GOARM, if we ever start building GOARM variants
	)
}

// genSPK generates an SPK, which is a nested tarball.
// The outer tar file is uncompressed and contains the minimal
// metadata. The main contents are in the outer tar's "package.tgz"
// entry, which is gzip or xz compressed.
func genSPK(w io.Writer, param spkParams) error {
	privFile := fmt.Sprintf("privilege-dsm%d", param.dsm)
	if param.forPackageCenter {
		privFile += ".priv"
	}

	var innerPkgTgz bytes.Buffer
	extractedSize, err := genInnerPackageTgz(&innerPkgTgz, param)
	if err != nil {
		return err
	}

	return writeTar(w,
		file("INFO", getInfo(param, extractedSize)),
		file("PACKAGE_ICON.PNG", static("PACKAGE_ICON.PNG")),
		file("PACKAGE_ICON_256.PNG", static("PACKAGE_ICON_256.PNG")),
		file("Tailscale.sc", static("Tailscale.sc")),
		dir("conf/", param.createTime),
		file("conf/resource", static("resource")),
		file("conf/privilege", static(privFile)),
		file("package.tgz", memFile(innerPkgTgz.Bytes(), 0644, param.createTime)),
		dir("scripts/", param.createTime),
		file("scripts/start-stop-status", static("scripts/start-stop-status")),
		file("scripts/postupgrade", static("scripts/postupgrade")),
		file("scripts/preupgrade", static("scripts/preupgrade")),
	)
}

func genInnerPackageTgz(w io.Writer, param spkParams) (extractedSize int64, err error) {
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
		file("bin/tailscaled", bin("tailscaled", param)),
		file("bin/tailscale", bin("tailscale", param)),
		dir("conf/", param.createTime),
		file("conf/Tailscale.sc", static("Tailscale.sc")),
		file("conf/logrotate.conf", static("logrotate-dsm"+strconv.Itoa(param.dsm))),
		dir("ui/", param.createTime),
		file("ui/PACKAGE_ICON_256.PNG", static("PACKAGE_ICON_256.PNG")),
		file("ui/index.cgi", bin("tailscale", param)), // TODO: don't build it again, don't include it twice
		file("ui/config", static("config")),           // TODO: this has "1.8.3" hard-coded in it; why? what is it? bug?
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

// bin returns a fileOpener for either the "tailscale" or "tailscaled"
// baseProg binary, building it or downloading it from pkgs.tailscale.com as
// necessary.
func bin(baseProg string, param spkParams) fileOpener {
	if param.srcDir == "" {
		return errLater(errors.New("TODO: download binaries from pkgs.tailscale.com"))
	}
	return buildBin(baseProg, param)
}

// buildBin builds baseProg ("tailscale" or "tailscaled")
// and returns a fileOpener for it.
func buildBin(baseProg string, param spkParams) fileOpener {
	vars, err := getDistVars(param.srcDir)
	if err != nil {
		return errLater(err)
	}

	cmd := exec.Command("go",
		"install",
		"-ldflags", ("-X tailscale.com/version.Long=" + vars.Long + " " +
			"-X tailscale.com/version.Short=" + vars.Short + " " +
			"-X tailscale.com/version.GitCommit=" + vars.GitHash),
		"tailscale.com/cmd/"+baseProg)
	cmd.Dir = param.srcDir
	cmd.Env = param.goEnv()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errLater(fmt.Errorf("building %s: %v, %s", baseProg, err, out))
	}
	cmd = exec.Command("go", "list", "-f", "{{.Target}}", "tailscale.com/cmd/"+baseProg)
	cmd.Dir = param.srcDir
	cmd.Env = param.goEnv()
	out, err = cmd.Output()
	if err != nil {
		return errLater(fmt.Errorf("running go list: %v", err))
	}
	binName := strings.TrimSpace(string(out))
	data, err := os.ReadFile(binName)
	if err != nil {
		return errLater(err)
	}
	return memFile(data, 0755, param.createTime)
}

type DistVars struct {
	Minor   string // "1.21"
	Short   string // "1.21.17"
	Long    string // "1.21.17-tb4f817065"
	GitHash string // "b4f8170657cde2a3a21ffee46c9dd028e400fb0f"
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
			sp = &v.Minor
		case "VERSION_SHORT":
			sp = &v.Short
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
	return v, bs.Err()
}

func getShortVer(dir string) (ver string, err error) {
	vars, err := getDistVars(dir)
	return vars.Short, err
}

// getInfo returns a fileOpener for the top-level INFO file.
// See genInfo.
func getInfo(param spkParams, extractedSize int64) fileOpener {
	data, err := genInfo(param, extractedSize)
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
func genInfo(param spkParams, extractedSize int64) ([]byte, error) {
	var buf bytes.Buffer
	add := func(k, v string) {
		fmt.Fprintf(&buf, "%s=%q\n", k, v)
	}
	add("package", "Tailscale")
	add("version", param.spkVersion())
	add("arch", param.synoArch)
	add("description", "Connect all your devices using WireGuard, without the hassle.")
	add("displayname", "Tailscale")
	add("maintainer", "Tailscale, Inc.")
	add("maintainer_url", "https://github.com/tailscale/tailscale-synology")
	add("create_time", param.createTime.Format("20060102-15:04:05"))
	add("dsmuidir", "ui")
	add("dsmappname", "SYNO.SDS.Tailscale")
	add("startstop_restart_services", "nginx")
	switch param.dsm {
	case 6:
		add("os_min_ver", "6.0.1-7445")
		add("os_max_ver", "7.0-40000")
	case 7:
		add("os_min_ver", "7.0-40000")
		add("os_max_ver", "")
	default:
		return nil, fmt.Errorf("unsupported DSM version '%v'", param.dsm)
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
