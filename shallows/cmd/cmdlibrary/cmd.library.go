package cmdlibrary

type cmdImports struct {
	Directory importDirectory `cmd:"" help:"import files from a directory into the library"`
}

type Commands struct {
	Import  cmdImports `cmd:"" help:"import media using various strategies"`
	Publist cmdPublish `cmd:"" help:"publish a library content"`
}
