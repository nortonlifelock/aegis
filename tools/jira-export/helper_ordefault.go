package main

import "time"

func tord(in *time.Time) (out time.Time) {
	if in != nil {
		out = *in
	}

	return out
}

func sord(in *string) (out string) {
	if in != nil {
		out = *in
	}

	return out
}

func ford(in *float32) (out float32) {
	if in != nil {
		out = *in
	}

	return out
}

func iord(in *int) (out int) {
	if in != nil {
		out = *in
	}

	return out
}
