package handler

import (
	"fmt"
	"goRent/internal/model"
)

// quickSort takes in an array and two ints as arguments and helps to sort the array passed in, in place.
func quickSort(array []int, left, right int) {

	if left < right {
		partitionIndex := partition(array, left, right)
		quickSort(array, left, partitionIndex-1)
		quickSort(array, partitionIndex+1, right)
	}

	fmt.Println("reached the end of quickSort")
}

func partition(array []int, left, right int) int {

	var partitionIndex int = left
	var pivotElement int = right

	for j := left; j < right; j++ {

		if array[j] < array[pivotElement] {

			swap(array, partitionIndex, j)
			partitionIndex++
		}
	}
	swap(array, partitionIndex, pivotElement)
	return partitionIndex
}

func swap(array []int, i1, i2 int) {
	temp := array[i1]
	array[i1] = array[i2]
	array[i2] = temp

}

// func mergeSort(a []model.SearchTrends) []model.SearchTrends {

// 	if len(a) < 2 {
// 		return a
// 	}

// 	left := 0
// 	right := len(a) - 1
// 	left_right := float64((left + right) / 2)
// 	middle := math.Floor(left_right)

// 	leftArray := a[:int(middle)]
// 	rightArray := a[int(middle):]

// 	return merge(mergeSort(leftArray), mergeSort(rightArray))
// }

// func merge(leftArray, rightArray []model.SearchTrends) []model.SearchTrends {

// 	leftIndex := 0
// 	rightIndex := 0

// 	var result []model.SearchTrends

// 	for leftIndex < len(leftArray) && rightIndex < len(rightArray) {

// 		if leftArray[leftIndex].Count < rightArray[rightIndex].Count {

// 			result = append(result, leftArray[leftIndex])
// 			leftIndex++
// 		} else {

// 			result = append(result, rightArray[rightIndex])
// 			rightIndex++
// 		}

// 	}

// 	for ; leftIndex < len(leftArray); leftIndex++ {
// 		result = append(result, leftArray[leftIndex])
// 	}

// 	for ; leftIndex < len(rightArray); leftIndex++ {
// 		result = append(result, rightArray[rightIndex])
// 	}

// 	return result

// }

// quickSortCategory takes in an array of model.SearchTrends and two ints as arguments and helps to sort the array passed in, in place.
func quickSortCategory(array []model.SearchTrends, left, right int) {

	if left < right {
		partitionIndex := partitionCategory(array, left, right)
		quickSortCategory(array, left, partitionIndex-1)
		quickSortCategory(array, partitionIndex+1, right)
	}

}

// partitionCategory works hand in hand with func quickSortCategory
func partitionCategory(array []model.SearchTrends, left, right int) int {

	var partitionIndex int = left
	var pivotElement int = right

	for j := left; j < right; j++ {

		if array[j].Count > array[pivotElement].Count {

			swapCategory(array, partitionIndex, j)
			partitionIndex++
		}
	}
	swapCategory(array, partitionIndex, pivotElement)
	return partitionIndex
}

// swapCategory works hand in hand with func partitionCategory
func swapCategory(array []model.SearchTrends, i1, i2 int) {
	temp := array[i1]
	array[i1] = array[i2]
	array[i2] = temp

}

// SortArrayCategory is the function to call when you want to call an array sorted. It takes in an array of type model.SearchTrends as an argument
func SortArrayCategory(a []model.SearchTrends) {

	left := 0
	right := len(a) - 1
	quickSortCategory(a, left, right)
}
