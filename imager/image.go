package imager

type Image struct {
	RootFsDir string  // root filesystem containing the files to bundle into the OS image
	Output string  // output directory/filename
	Format string  // output format
	// Volume string
	// Compression string

	BuildDir string  // the build directory use to store the temporary files we need to build the image
}
