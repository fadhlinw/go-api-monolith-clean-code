package utils

import (
	"errors"
	"fmt"
	"strconv"
)

// QRISData holds the decoded QRIS information.
type QRISData struct {
	Tag     string
	Length  int
	Value   string
	SubTags map[string]QRISData
}

// DecodeQRIS decodes the static QRIS string into a structured format.
func DecodeQRIS(qris string) (map[string]QRISData, error) {
	result := make(map[string]QRISData)
	index := 0

	for index < len(qris) {
		// Ensure we have at least 4 characters for tag and length.
		if len(qris)-index < 4 {
			return nil, fmt.Errorf("invalid QRIS string format: insufficient length at index %d", index)
		}

		// Parse Tag
		tag := qris[index : index+2]
		index += 2

		// Parse Length
		lengthStr := qris[index : index+2]
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid length at tag %s: %v", tag, err)
		}
		index += 2

		// Parse Value
		if len(qris)-index < length {
			return nil, fmt.Errorf("value exceeds string length at tag %s", tag)
		}
		value := qris[index : index+length]
		index += length

		// Decode sub-tags if applicable
		subTags := make(map[string]QRISData)
		if tag == "26" {
			subTags, err = decodeSubTags(value)
			if err != nil {
				return nil, fmt.Errorf("failed to decode sub-tags at tag %s: %v", tag, err)
			}
		}

		// Add the parsed data to the result map
		result[tag] = QRISData{
			Tag:     tag,
			Length:  length,
			Value:   value,
			SubTags: subTags,
		}
	}

	return result, nil
}

// decodeSubTags decodes nested tags within a value.
func decodeSubTags(value string) (map[string]QRISData, error) {
	subTags := make(map[string]QRISData)
	index := 0

	for index < len(value) {
		// Ensure we have at least 4 characters for sub-tag and length.
		if len(value)-index < 4 {
			return nil, errors.New("invalid sub-tag format: insufficient length")
		}

		// Parse Sub-Tag
		subTag := value[index : index+2]
		index += 2

		// Parse Length
		lengthStr := value[index : index+2]
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid sub-tag length for sub-tag %s: %v", subTag, err)
		}
		index += 2

		// Parse Value
		if len(value)-index < length {
			return nil, fmt.Errorf("sub-tag value exceeds string length for sub-tag %s", subTag)
		}
		subValue := value[index : index+length]
		index += length

		// Add the parsed sub-tag to the map
		subTags[subTag] = QRISData{
			Tag:    subTag,
			Length: length,
			Value:  subValue,
		}
	}

	return subTags, nil
}
