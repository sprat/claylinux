package imager

/*
func CreateRawDisk(printf func(string, ...any), path string, diskSize int64) error {
    printf("creating raw disk of size %s", humanize.Bytes(uint64(diskSize)))

    f, err := os.Create(path)
    if err != nil {
        return fmt.Errorf("failed to create raw disk: %w", err)
    }

    defer f.Close() //nolint:errcheck

    if err = f.Truncate(diskSize); err != nil {
        return fmt.Errorf("failed to create raw disk: %w", err)
    }

    if err = syscall.Fallocate(int(f.Fd()), 0, 0, diskSize); err != nil {
        fmt.Fprintf(os.Stderr, "WARNING: failed to preallocate disk space for %q (size %d): %s", path, diskSize, err)
    }

    return f.Close()
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
    run("mcopy", "-i", espFile, efiFile, "::/EFI/boot/boot"+efiEfiArch+".efi")
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
