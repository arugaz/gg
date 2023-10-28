// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import "math"

// Matrix represents a 2D transformation matrix with six components: XX, YX, XY, YY, X0, and Y0.
// This matrix is used to transform points, vectors, or perform various geometric transformations.
type Matrix struct {
	XX, YX, XY, YY, X0, Y0 float64
}

// Identity returns the identity matrix, which has no effect on transformations.
// The identity matrix is used as a starting point for geometric transformations.
func Identity() Matrix {
	return Matrix{
		1, 0,
		0, 1,
		0, 0,
	}
}

// Translate returns a translation matrix that moves points by the specified
// horizontal and vertical offsets (x, y). This matrix is used to shift objects
// in a 2D space.
func Translate(x, y float64) Matrix {
	return Matrix{
		1, 0,
		0, 1,
		x, y,
	}
}

// Scale returns a scaling matrix that scales points by the specified horizontal
// and vertical factors (x, y). This matrix is used to resize objects in a 2D space.
func Scale(x, y float64) Matrix {
	return Matrix{
		x, 0,
		0, y,
		0, 0,
	}
}

// Rotate returns a rotation matrix that rotates points by the specified angle
// (in radians) in a counterclockwise direction. This matrix is used to rotate objects
// in a 2D space.
func Rotate(angle float64) Matrix {
	c := math.Cos(angle)
	s := math.Sin(angle)

	return Matrix{
		c, s,
		-s, c,
		0, 0,
	}
}

// Shear returns a shear matrix that shears points in the horizontal (x) and
// vertical (y) directions. This matrix is used to skew or slant objects in a 2D space.
func Shear(x, y float64) Matrix {
	return Matrix{
		1, y,
		x, 1,
		0, 0,
	}
}

// Multiply multiplies two matrices, a and b, and returns a new matrix that represents
// the combined transformation. This operation is used to apply a sequence of matrix
// transformations to points or objects in a 2D space.
func (a Matrix) Multiply(b Matrix) Matrix {
	return Matrix{
		a.XX*b.XX + a.YX*b.XY,
		a.XX*b.YX + a.YX*b.YY,
		a.XY*b.XX + a.YY*b.XY,
		a.XY*b.YX + a.YY*b.YY,
		a.X0*b.XX + a.Y0*b.XY + b.X0,
		a.X0*b.YX + a.Y0*b.YY + b.Y0,
	}
}

// TransformVector applies the matrix transformation represented by 'a' to a 2D vector (x, y).
// It returns the transformed vector (tx, ty) after applying the matrix transformation.
func (a Matrix) TransformVector(x, y float64) (tx, ty float64) {
	tx = a.XX*x + a.XY*y
	ty = a.YX*x + a.YY*y

	return tx, ty
}

// TransformPoint applies the matrix transformation represented by 'a' to a 2D point (x, y).
// It returns the transformed point (tx, ty) after applying the matrix transformation.
func (a Matrix) TransformPoint(x, y float64) (tx, ty float64) {
	tx = a.XX*x + a.XY*y + a.X0
	ty = a.YX*x + a.YY*y + a.Y0

	return tx, ty
}

// Translate returns a new matrix resulting from applying a 2D translation by (x, y) to matrix 'a'.
func (a Matrix) Translate(x, y float64) Matrix {
	return Translate(x, y).Multiply(a)
}

// Scale returns a new matrix resulting from applying a 2D scaling transformation to matrix 'a'.
func (a Matrix) Scale(x, y float64) Matrix {
	return Scale(x, y).Multiply(a)
}

// Rotate returns a new matrix resulting from applying a 2D rotation transformation by the given angle to matrix 'a'.
func (a Matrix) Rotate(angle float64) Matrix {
	return Rotate(angle).Multiply(a)
}

// Shear returns a new matrix resulting from applying a 2D shear transformation by the given factors to matrix 'a'.
func (a Matrix) Shear(x, y float64) Matrix {
	return Shear(x, y).Multiply(a)
}
