xlattice_go
===========

An implementation of [XLattice](http://xlattice.sourceforge.net)
for the Go language.  XLattice is a communications library 
for peer-to-peer networks.  More extensive (although somewhat
dated) information on XLattice is available as the 
[XLattice website](http://www.xlattice.org).

This version of xlattice-go includes a Go implementation of *rnglib*
for use in testing and an implementation of u256x256, a system for
storing files by their content keys.

## u256x256

Rather than storing data files in a hierarchical directory structure
where both directories and data files are given string names, u256x256
stores files named by their content keys.  The keys are generated by
either SHA1 or SHA3 (Keccak-256).  Storage by content keys has several 
advantages.  For one, it is trivial to determine whether a file is corrupt:
you simply recalculate the hash.  In a distributed storage system, files
are requested by key.  All machines participating in the retrieval can
check file integrity as the file is passing through and drop and
re-request if the hash doesn't match the content key.

u256x256 is optimized for storing very large numbers of files.  The
first byte of the content key determines which top-level directory
the file goes in; the second byte determines its lower-level 
directory.  So if a file's content hash is abcdef...1234, then it
will be stored in ab/cd/ef...1234.  There are 256 top-level 
directories and 256 subdirectories below each of these, so 256x256 = 
65536 lower-level directories.  

	// Determine the SHA1 or SHA3 content hash of an arbitrary file
	func FileSHA1(path string) (hash string, err error)
	func FileSHA3(path string) (hash string, err error)
	
	// Create a u256x256 directory structure
	func New(path string) *U256x256
	
	// Attributes for files in U, a u256x256 directory tree
	func (u *U256x256) Exists(key string) bool
	func (u *U256x256) FileLen(key string) (length int64, err error)
	func (u *U256x256) GetPathForKey(key string) string
	
	// Copy a data file and add the copy to U using an SHA1 key.  If the
	// key doesn't match, the operation fails.
	func (u *U256x256) CopyAndPut1(path, key string) (
	    written int64, hash string, err error)
	// Retrieve a file by its SHA1 key.
	func (u *U256x256) GetData1(key string) (
	    data []byte, err error)
	// Insert a data file into U; the original is lost.
	func (u *U256x256) Put1(inFile, key string) (
	    length int64, hash string, err error)
	// Write a buffer into U, storing it by its SHA1 key.
	func (u *U256x256) PutData1(data []byte, key string) (
	    length int64, hash string, err error)
	
	// Similar functions using the SHA3 (Keccak-256) hash function.
	func (u *U256x256) CopyAndPut3(path, key string) (
	    written int64, hash string, err error)
	func (u *U256x256) GetData3(key string) (
	    data []byte, err error){ var path string
	func (u *U256x256) Put3(inFile, key string) (
	    length int64, hash string, err error)
	func (u *U256x256) PutData3(data []byte, key string) (
	    length int64, hash string, err error)

## rnglib

A Go random number generator especially useful for generating 
quasi-random strings, files, and directories.

This package contains two classes, SimpleRNG and SystemRNG.
Both of these subclass Go's rand package and so rand's functions 
can be called through either of rnglib's two subclasses using 
exactly the same syntax.

SimpleRNG is the faster of the two.  It uses the Mersenne
Twister and so is completely predictable with a very long period.
It is suitable where speed and predictability are both important.
You can be certain that if you provide the same seed, then
the sequence of numbers generated will be exactly the same.  This is
often important for debugging.  

SimpleRNG is also very fast.  It is approximately 30% faster 
than Go's rand package in our tests.

Although the Mersenne Twister is absolutely predictable, it also has 
extremely good statistical properties.  Tests using standard
packages such as [Diehard](http://en.wikipedia.org/wiki/Diehard_tests)
and [Dieharder](http://www.phy.duke.edu/~rgb/General/dieharder.php) rank
the Mersenned Twister very highly.  The Twister is quite widely used 
(for example, in the Python libraries) because of this.

SystemRNG is a secure random number generator, meaning that it is
extremely difficult or impossible to predict the next number 
generated in a sequence.  It is based on the system's /dev/urandom.
This relies upon entropy accumulated by the system; when 
this is insufficient it will revert to using a very secure 
programmable random number generator.  SystemRNG is about 35x
slower than SimpleRNG in our test runs.

For normal use in testing, SimpleRNG is preferred.  When there is
a requirement for greater security, as in the generation of passwords 
and cryptographic keys with reasonable strength, SystemRNG is recommended.  

In addition to the functions available from Go's random package,
rnglib also provides

    NextBoolean()
    NextByte(max=256)
    NextBytes(buffer)
    NextInt32(n uint32)
    NextInt64(n uint64)
    NextFloat32()
    NextFloat64()
    NextFileName(maxLen int)
    NextDataFile(dirName string, maxLen, minLen int)
    NextDataDir(pathToDir string, depth, width, maxLen, minLen int)

Given a buffer of arbitrary length, nextBytes() will fill it with random
bytes.

File names generated by nextFileName() are at least one character long.  
The first letter must be alphabetic (including the underscore) 
whereas other characters may also be numeric or either of dot and dash
("." and "-").

The function nextDataFile() creates a file with a random name in the
directory indicated.  The file length will vary from minLen (which 
defaults to zero) up to but excluding maxLen bytes.

The nextDataDir() function creates a directory structure that is
'depth' deep and 'width' wide, where the latter means that there 
will be that many files in each directory and subdirectory created
(where is file is either a data file or a directory).  Depth and
width must both be at least 1.  maxLen and minLen are used as in
nextDataFile().

## NOTICE

Most of the above projects have been split off as separate Github
repositories.  The code here will no longer be maintained and will at
some point be deleted.  The new repositories are:

+ <https://github.com/jddixon/rnglib_go>
+ <https://github.com/jddixon/xlU_go>
+ <https://github.com/jddixon/xlUtil_go>
+ <https://github.com/jddixon/xlCrypto_go>
+ <https://github.com/jddixon/xlTransport_go>
+ <https://github.com/jddixon/xlProtocol_go>
+ <https://github.com/jddixon/xlOverlay_go>
+ <https://github.com/jddixon/xlNodeID_go>
+ <https://github.com/jddixon/xlNode_go>
+ <https://github.com/jddixon/xlReg_go>
