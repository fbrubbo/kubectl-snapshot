package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// String2MilliCPU converts String to Milli CPU
func String2MilliCPU(cpu string) int {
	if strings.Contains(cpu, "m") {
		str := strings.ReplaceAll(cpu, "m", "")
		milli, _ := strconv.Atoi(str)
		return milli
	}
	cpuf, _ := strconv.ParseFloat(cpu, 64)
	milli := (int)(cpuf * 1000)
	return milli
}

// String2MiMemory converts String to Memory Mi
func String2MiMemory(men string) int {
	reg, _ := regexp.Compile(`(\d*)(.*)`)
	groups := reg.FindStringSubmatch(men)
	memory, _ := strconv.Atoi(groups[1])
	suffix := groups[2]

	switch suffix {
	case "G":
		// http://extraconversion.com/data-storage-conversion-table/gigabytes-to-mebibytes.html
		return int(math.Round(float64(memory) * 953.67431640625))
	case "Gi":
		// http://extraconversion.com/data-storage-conversion-table/gibibytes-to-mebibytes.html
		return memory * 1024
	case "M":
		// http://extraconversion.com/data-storage-conversion-table/megabytes-to-mebibytes.html
		return int(math.Round(float64(memory) * 0.9537))
	case "Mi":
		return memory
	case "Ki":
		// http://extraconversion.com/data-storage-conversion-table/mebibytes-to-kibibytes.html
		return int(math.Round(float64(memory) / 1024))
	default:
		// http://extraconversion.com/data-storage-conversion-table/bytes-to-mebibytes.html
		return int(math.Round(float64(memory) * 9.53674E-7))
	}

	/*
		TODO:

		Limits and requests for memory are measured in bytes.
		You can express memory as a plain integer or as a fixed-point integer using one of these suffixes: E, P, T, G, M, K.
		You can also use the power-of-two equivalents: Ei, Pi, Ti, Gi, Mi, Ki. For example, the following represent roughly the same value:

		128974848, 129e6, 129M, 123Mi
	*/
}
