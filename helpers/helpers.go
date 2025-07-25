// Copyright (c) 2025, Zainal Abidin
// SPDX-License-Identifier: BSD-2-Clause

package helpers

func RoundToTwoDecimals(val float64) float64 {
	return float64(int(val*100)) / 100
}