package wordle

import (
  "os"
  "fmt"
  "bytes"
)

func isValid(target [5]byte) (bool, error) {
  // Open the file containing the word list
  file, err := os.Open("data/wordle/valid-words.txt")
	if err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}
  defer file.Close()

  // Initialize binary search bounds
  low, high := 1, 12972
  for low <= high {
    mid := (high + low) / 2

    // Read the word at the midpoint of the current range
    word, err := readWord(file, mid)
	  if err != nil {
		  return false, fmt.Errorf("error reading file: %w", err)
	  }

    // Compare the current word with the target word
    result := bytes.Compare(word, target[:])

    // Return true if the word matches the target
    if result == 0 {
      return true, nil

    // If the current word is lexicographically smaller, ignore left half
    } else if result < 0 {
      low = mid + 1;

    // If the current word is lexicographically larger, ignore right half
    } else {
      high = mid - 1;
    }
  }

  return false, nil
}

func readWord(file *os.File, line int) ([]byte, error) {
  // Calculate the byte offset for the specified line in the file
  offset := int64(line - 1) * 6 // 6 bytes per word (5 letters + newline)

  // Seek to the calculated offset
  _, err := file.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking to offset: %w", err)
	}

  // Read the next 5 bytes (the word itself)
  word := make([]byte, 5)
  _, err = file.Read(word)
	if err != nil {
		return nil, fmt.Errorf("error reading word from file: %w", err)
	}

  return word, nil
}
