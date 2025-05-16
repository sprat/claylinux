package main

import (
	"log"
	"os"

	"github.com/sprat/claylinux/imager/efi"
)

/*
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("could not stat file %s: %v", path, err)
	}
	return info.Size(), nil
}

func inMib(n int64) int64 {
	return (n + (1 << 20) - 1) >> 20
}

func FileSizeMib(path string) (int64, error) {
	size, err := getFileSize(path)
	return inMib(size), err
}

func align(val, multiple int64) int64 {
	return ((val + multiple - 1) / multiple) * multiple
}

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		die(fmt.Sprintf("command failed: %s %v: %v", name, args, err))
	}
}

func runOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		die(fmt.Sprintf("command failed: %s %v: %v", name, args, err))
	}
	return strings.TrimSpace(string(out))
}

// --- Core build steps ---
func compress(filename string) {
	switch compression {
	case "none":
		// do nothing
	case "gz":
		run("pigz", "-9", filename)
		run("mv", filename+".gz", filename)
	case "xz":
		run("xz", "-C", "crc32", "-9", "-T0", filename)
		run("mv", filename+".xz", filename)
	case "zstd":
		run("zstd", "-19", "-T0", "--rm", filename)
		run("mv", filename+".zstd", filename)
	default:
		die("invalid compression scheme: " + compression)
	}
}

// --- Build steps ---

func buildInitramfs() {
	os.Mkdir("initramfs_files", 0755)
	run("cp", "/usr/share/claylinux/init", "initramfs_files")

	if fileExists("/system/etc/hosts.target") {
		os.MkdirAll("initramfs_files/etc", 0755)
		run("cp", "/system/etc/hosts.target", "initramfs_files/etc/hosts")
	}
	if fileExists("/system/etc/resolv.conf.target") {
		os.MkdirAll("initramfs_files/etc", 0755)
		run("cp", "/system/etc/resolv.conf.target", "initramfs_files/etc/resolv.conf")
	}
	// cpio for initramfs_files
	run("sh", "-c", `find initramfs_files -mindepth 1 -printf '%P\0' | cpio --quiet -o0H newc -D initramfs_files -F initramfs.img`)
	// Add system files except boot, hosts.target, resolv.conf.target
	run("sh", "-c", `find /system -path /system/boot -prune -o ! -path /system/init ! -path /system/etc/hosts.target ! -path /system/etc/resolv.conf.target -mindepth 1 -printf '%P\0' | cpio --quiet -o0AH newc -D /system -F initramfs.img`)
	compress("initramfs.img")

	ucode := runOutput("find", "/system/boot/", "-name", "*-ucode.img")
	imgs := "initramfs.img"
	if ucode != "" {
		imgs = ucode + " " + imgs
	}
	run("sh", "-c", fmt.Sprintf("cat %s >initramfs", imgs))
	run("find", ".", "!", "-name", "initramfs", "-delete")
}

// getInitialOffset computes offset+size from objdump output.
func getInitialOffset(efiStub string) (int64, error) {
	cmd := exec.Command("objdump", "-h", "-w", efiStub)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	var lastFields []string
	for scanner.Scan() {
		line := scanner.Text()
		// Sections have lines that start with space+number or name, skip headers
		fields := strings.Fields(line)
		if len(fields) >= 5 && strings.HasPrefix(fields[1], "0x") {
			lastFields = fields
		}
	}
	if len(lastFields) < 5 {
		return 0, fmt.Errorf("failed to parse objdump output")
	}
	offset, err := strconv.ParseInt(lastFields[4], 0, 64)
	if err != nil {
		return 0, err
	}
	size, err := strconv.ParseInt(lastFields[3], 0, 64)
	if err != nil {
		return 0, err
	}
	return int64(offset + size), nil
}

func buildUKI() {
	const alignment = 4096
	efiStub := "path/to/efi_stub" // Set your input stub
	efiFile := "path/to/efi_file" // Set your output file

	// For example, sections := [][2]string{{".foo", "foo.bin"}, {".bar", "bar.bin"}}
	sections := [][2]string{
		{".section1", "file1.bin"},
		{".section2", "file2.bin"},
		// Add more as needed
	}

	// Step 1: Get initial offset from objdump and align it
	offset, err := getInitialOffset(efiStub)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting initial offset: %v\n", err)
		os.Exit(1)
	}
	offset = align(offset, alignment)

	// Step 2: Prepare objcopy arguments
	var args []string
	for _, pair := range sections {
		section, file := pair[0], pair[1]
		args = append(args, "--add-section", fmt.Sprintf("%s=%s", section, file))
		args = append(args, "--change-section-vma", fmt.Sprintf("%s=0x%X", section, offset))

		size, err := getFileSize(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting size of %s: %v\n", file, err)
			os.Exit(1)
		}
		size = align(size, alignment)
		offset += size
	}

	// Step 3: Run objcopy
	args = append(args, efiStub, efiFile)
	cmd := exec.Command("objcopy", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("Running: objcopy %s\n", strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "objcopy failed: %v\n", err)
		os.Exit(1)
	}

}

func buildEFI() {
	os.Chdir(buildDir)
	fmt.Println("Building the EFI executable")
	buildInitramfs()
	size := getFileSizeMib("initramfs")
	fmt.Printf("The size of the initramfs is: %d MiB\n", size)
	kernel := runOutput("find", "/system/boot", "-name", "vmlinu*", "-print")
	run("sh", "-c", "tr '\\n' ' ' </system/boot/cmdline >cmdline")
	run("sh", "-c", "basename /system/lib/modules/* >kernel-release")

	// Compose the UKI section file
	ukiStanza := fmt.Sprintf(
		".osrel /system/etc/os-release\n.uname kernel-release\n.cmdline cmdline\n.initrd initramfs\n.linux %s\n", kernel)
	cmd := exec.Command("build_uki")
	cmd.Stdin = strings.NewReader(ukiStanza)
	if err := cmd.Run(); err != nil {
		die("build_uki failed: " + err.Error())
	}

	run("find", ".", "!", "-name", "*.efi", "-delete")
}

func generateEFI() {
	fmt.Println("Copying the OS files to the output directory")
	run("mv", efiFile, output+".efi")
}

func generateESP() {
	fmt.Println("Generating the EFI System Partition (ESP)")
	size := getFileSize(efiFile)
	size = inMib(size * 102 / 100)
	fmt.Printf("The size of the ESP is: %d MiB\n", size)
	run("mkfs.vfat", "-n", volume, "-F", "32", "-C", espFile, strconv.FormatInt(size<<10, 10))
	run("mmd", "-i", espFile, "::/EFI")
	run("mmd", "-i", espFile, "::/EFI/boot")
	run("mcopy", "-i", espFile, efiFile, "::/EFI/boot/boot"+efiArch+".efi")
	run("rm", efiFile)
}

func generateISO() {
	generateESP()
	fmt.Println("Generating the ISO image")
	run("xorrisofs", "-e", filepath.Base(espFile),
		"-no-emul-boot", "-joliet", "-full-iso9660-filenames", "-rational-rock",
		"-sysid", "LINUX", "-volid", volume,
		"-output", output+".iso", espFile)
	run("rm", espFile)
}

func generateRaw() {
	generateESP()
	fmt.Println("Generating the disk image")
	diskFile := output + ".img"
	espMiB := getFileSizeMib(espFile)
	diskSize := espMiB + 2
	fmt.Printf("The size of the disk image is: %d MiB\n", diskSize)
	run("truncate", "-s", fmt.Sprintf("%dM", diskSize), diskFile)
	sfdiskCmd := fmt.Sprintf("label: gpt\nfirst-lba: 34\nstart=1MiB size=%dMiB name=\"EFI system partition\" type=uefi\n", espMiB)
	cmd := exec.Command("sfdisk", "--quiet", diskFile)
	cmd.Stdin = strings.NewReader(sfdiskCmd)
	if err := cmd.Run(); err != nil {
		die("sfdisk failed: " + err.Error())
	}
	run("dd", "if="+espFile, "of="+diskFile, "bs=1M", "seek=1", "conv=notrunc", "status=none")
	run("rm", espFile)
}

func convertImage(format string) {
	generateRaw()
	fmt.Printf("Converting the disk image to %s format\n", format)
	run("qemu-img", "convert", "-f", "raw", "-O", format, output+".img", output+"."+format)
	run("rm", output+".img")
}
*/

// --- Option parsing, main ---
func main() {
	/*
	output = "/out/claylinux"
	format = "raw"
	volume = "CLAYLINUX"
	compression = "gz"

	flag.StringVar(&format, "format", format, "Output format (efi, iso, raw, qcow2, vmdk, vhdx, vdi)")
	flag.StringVar(&output, "output", output, "Output image path/name, without extension")
	flag.StringVar(&volume, "volume", volume, "Volume label for the boot partition")
	flag.StringVar(&compression, "compression", compression, "Compression format for initramfs: none|gz|xz|zstd")
	flag.Parse()

	if _, err := os.Stat("/system"); err != nil {
		die("the /system directory does not exist, please copy/mount your root filesystem here")
	}
	*/

	log.Printf("EFI architecture: %s", efi.Suffix)
	log.Printf("Default EFI name: %s", efi.DefaultName)
	log.Printf("EFI stub: %s", efi.Stub)

	buildDir, err := os.MkdirTemp("", "claylinux-imager-")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Temporary build directory: %s", buildDir)
	defer os.RemoveAll(buildDir)

	/*
	efiFile = filepath.Join(buildDir, "claylinux.efi")
	espFile = filepath.Join(buildDir, "claylinux.esp")

	os.Chdir(buildDir)
	buildEFI()
	os.Chdir("/")

	os.MkdirAll(filepath.Dir(output), 0755)
	switch format {
	case "efi":
		generateEFI()
	case "iso":
		generateISO()
	case "raw":
		generateRaw()
	case "qcow2", "vmdk", "vhdx", "vdi":
		convertImage(format)
	default:
		die("invalid format: " + format)
	}
	*/
}
