package common

import (
	"crypto/md5"
	"io"
)

// ComputeMD5 ComputeMD5
func ComputeMD5(src io.Reader) (digest []byte, err error) {
	hash := md5.New()
	buffer := make([]byte, 32*1024)
	_, err = io.CopyBuffer(hash, src, buffer)
	if err != nil {
		return nil, err
	}
	return hash.Sum(make([]byte, 0)), nil
}

// CopyAndComputeMD5 CopyAndComputeMD5
func CopyAndComputeMD5(
	dst io.Writer,
	src io.Reader,
) (digest []byte, written int64, err error) {
	digest = nil
	written = 0
	err = nil
	hash := md5.New()
	buffer := make([]byte, 32*1024)
	for {
		srcRead, srcErr := src.Read(buffer)
		if srcRead > 0 {
			dstWritten, dstErr := dst.Write(buffer[0:srcRead])
			hansWritten, hashErr := hash.Write(buffer[0:srcRead])
			if dstWritten > 0 {
				written += int64(dstWritten)
			}
			if dstErr != nil {
				err = dstErr
				break
			}
			if hashErr != nil {
				err = hashErr
				break
			}
			if srcRead != dstWritten || srcRead != hansWritten {
				err = io.ErrShortWrite
				break
			}
		}
		if srcErr != nil {
			if srcErr != io.EOF {
				err = srcErr
			}
			break
		}
	}
	if err == nil {
		digest = hash.Sum(make([]byte, 0))
	}
	return digest, written, err
}
