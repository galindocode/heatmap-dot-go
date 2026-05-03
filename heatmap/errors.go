package heatmap

import "errors"

// Common errors returned by the heatmap package
var (
	// ErrInvalidSize is returned when width or height is invalid (<=0)
	ErrInvalidSize = errors.New("heatmap: invalid size, width and height must be > 0")

	// ErrNoData is returned when trying to generate a heatmap with no data points
	ErrNoData = errors.New("heatmap: no data points provided")

	// ErrInvalidMaxValue is returned when MaxValue is set to <= 0
	ErrInvalidMaxValue = errors.New("heatmap: MaxValue must be > 0")

	// ErrInvalidPointSize is returned when PointSize is <= 0
	ErrInvalidPointSize = errors.New("heatmap: PointSize must be > 0")

	// ErrInvalidColor is returned when a color hex string is invalid
	ErrInvalidColor = errors.New("heatmap: invalid color format")

	// ErrInvalidImage is returned when a base image is nil or invalid
	ErrInvalidImage = errors.New("heatmap: invalid base image")

	// ErrSizeMismatch is returned when overlay image size doesn't match heatmap size
	ErrSizeMismatch = errors.New("heatmap: image size does not match heatmap dimensions")
)
