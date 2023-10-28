package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Point struct {
	row, col int
}

func main() {
	path := strings.Join(os.Args[1:], " ")
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	str := string(bytes)
	lines := strings.Split(str, "\n")
	sudoku := [][]int{}
	for i := 0; i < len(lines); i++ {
		result := []int{}
		nums := strings.Split(strings.Join(strings.Split(lines[i], "|"), " "), " ")
		for j := 0; j < len(nums); j++ {
			num, err := strconv.Atoi(nums[j])
			if err == nil {
				result = append(result, num)
			}
		}
		sudoku = append(sudoku, result)
	}

	result := solveSudoku(sudoku)
	resultText := ""
	for row := 0; row < len(result); row++ {
		line := "|"
		for col := 0; col < len(result); col++ {
			line += strconv.Itoa(result[row][col])
			if col > 0 && (col+1)%3 == 0 {
				line += "|"
			} else {
				line += " "
			}
		}
		resultText += line + "\n"
	}
	ext := filepath.Ext(path)
	file, err := os.Create(strings.TrimSuffix(path, ext) + "_solved.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(resultText[:len(resultText)-1])
	if err != nil {
		panic(err)
	}
}

func solveSudoku(puzzle [][]int) [][]int {
	puzzleCopy := make([][]int, len(puzzle))
	copy(puzzleCopy, puzzle)
	activeNumbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rectCentres := getRectCentres(len(puzzleCopy), len(puzzleCopy[0]))
	columns := getColumns(puzzleCopy)
	points := map[string][]int{}
	for row := 0; row < len(puzzleCopy); row++ {
		for col := 0; col < len(puzzleCopy); col++ {
			number := puzzleCopy[row][col]
			if number > 0 {
				continue
			}
			activeRectCenter := filter(rectCentres, func(c Point) bool {
				return pointInRect(Point{row: row, col: col}, c)
			})[0]
			activeRectNumbers := getRectNumbers(activeRectCenter, puzzleCopy)
			numberVariants := filter(activeNumbers, func(n int) bool {
				checkFunc := func(el int) bool {
					return el == n
				}
				return (!some(columns[col], checkFunc) &&
					!some(puzzleCopy[row], checkFunc) &&
					!some(activeRectNumbers, checkFunc))
			})
			if len(numberVariants) > 0 {
				points[createKey(row, col)] = numberVariants
			}
		}
	}

	for len(points) > 0 {
		emptyPointsKeys := getKeys(points)
		if len(emptyPointsKeys) == 0 {
			break
		}

		for i := 0; i < len(emptyPointsKeys); i++ {
			key := emptyPointsKeys[i]
			numberVariants := points[key]
			if len(numberVariants) > 1 {
				continue
			}
			row, col := getRowAndColFromKey(key)
			settedNumber := numberVariants[0]
			puzzleCopy[row][col] = settedNumber
			delete(points, key)
			center := filter(rectCentres, func(el Point) bool {
				return pointInRect(Point{row: row, col: col}, el)
			})[0]
			compareFunc := func(e int) bool {
				return e != settedNumber
			}
			for r := 0; r < len(puzzleCopy); r++ {
				key := createKey(r, col)
				if _, ok := points[key]; ok {
					points[key] = filter(points[key], compareFunc)
				}
			}
			for c := 0; c < len(puzzleCopy); c++ {
				key := createKey(row, c)
				if _, ok := points[key]; ok {
					points[key] = filter(points[key], compareFunc)
				}
			}
			for row := center.row - 1; row <= center.row+1; row++ {
				for col := center.col - 1; col <= center.col+1; col++ {
					key := createKey(row, col)
					if _, ok := points[key]; ok {
						points[key] = filter(points[key], compareFunc)
					}
				}
			}
		}
	}

	return puzzleCopy
}

func getRectNumbers(center Point, puzzle [][]int) []int {
	result := []int{}
	for row := center.row - 1; row <= center.row+1; row++ {
		for col := center.col - 1; col <= center.col+1; col++ {
			result = append(result, puzzle[row][col])
		}
	}
	return result
}

func pointInRect(point Point, rectCenter Point) bool {
	return (rectCenter.row-1 <= point.row &&
		point.row <= rectCenter.row+1 &&
		rectCenter.col-1 <= point.col &&
		point.col <= rectCenter.col+1)
}

func getRectCentres(rowCount int, colCount int) []Point {
	size := 3
	result := []Point{}
	for row := 0; row < rowCount; row += size {
		for col := 0; col < colCount; col += size {
			result = append(result, Point{
				row: int(math.Floor(float64(row) + float64(size)/2)),
				col: int(math.Floor(float64(col) + float64(size)/2)),
			})
		}
	}
	return result
}

func getColumns(puzzle [][]int) [][]int {
	result := [][]int{}
	for i := 0; i < len(puzzle); i++ {
		result = append(result,
			mapArray(puzzle, func(el []int) int {
				return el[i]
			}),
		)
	}
	return result
}

func getKeys(mapper map[string][]int) []string {
	keys := make([]string, 0, len(mapper))
	for k := range mapper {
		keys = append(keys, k)
	}
	return keys
}

func createKey(row int, col int) string {
	return strconv.Itoa(row) + "-" + strconv.Itoa(col)
}

func getRowAndColFromKey(key string) (int, int) {
	rowAndCol := mapArray(strings.Split(key, "-"), func(el string) int {
		num, err := strconv.Atoi(el)
		if err != nil {
			panic(err)
		}
		return num
	})
	return rowAndCol[0], rowAndCol[1]
}

func filter[T any](arr []T, filterFunc func(T) bool) []T {
	result := []T{}
	for i := 0; i < len(arr); i++ {
		if filterFunc(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

func mapArray[T any, Tres any](arr []T, mapper func(T) Tres) []Tres {
	result := make([]Tres, len(arr))
	for i := 0; i < len(arr); i++ {
		result[i] = mapper(arr[i])
	}
	return result
}

func some[T any](arr []T, checkFunc func(T) bool) bool {
	for i := 0; i < len(arr); i++ {
		if checkFunc(arr[i]) {
			return true
		}
	}
	return false
}
