package main

import (
	"reflect"
	"testing"
)

func TestColor_String(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{
			name:  "Test Blue color",
			color: Blue,
			want:  "Blue",
		},
		{
			name:  "Test Red color",
			color: Red,
			want:  "Red",
		},
		{
			name:  "Test Green color",
			color: Green,
			want:  "Green",
		},
		{
			name:  "Test Yellow color",
			color: Yellow,
			want:  "Yellow",
		},
		{
			name:  "Test Pink color",
			color: Pink,
			want:  "Pink",
		},
		{
			name:  "Test White color",
			color: White,
			want:  "White",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.String(); got != tt.want {
				t.Errorf("Color.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcHB(t *testing.T) {
	type args struct {
		answer Line
		sub    Line
	}
	tests := []struct {
		name string
		args args
		want HitBlow
	}{
		{
			name: "4Hits/0Blows (3 Colors)",
			args: args{
				answer: Line{Red, Red, Blue, Green},
				sub:    Line{Red, Red, Blue, Green},
			},
			want: HitBlow{nhit: 4, nblow: 0},
		},
		{
			name: "1Hits/1Blows (3 Colors)",
			args: args{
				answer: Line{Red, Red, Blue, Green},
				sub:    Line{Blue, White, Blue, Red},
			},
			want: HitBlow{nhit: 1, nblow: 1},
		},
		{
			name: "0Hits/4Blows (4 Colors)",
			args: args{
				answer: Line{Red, White, Blue, Green},
				sub:    Line{White, Blue, Green, Red},
			},
			want: HitBlow{nhit: 0, nblow: 4},
		},
		{
			name: "0Hits/4Blows (2 Colors)",
			args: args{
				answer: Line{Red, Red, Yellow, Yellow},
				sub:    Line{Yellow, Yellow, Red, Red},
			},
			want: HitBlow{nhit: 0, nblow: 4},
		},
		{
			name: "0Hits/0Blows (2 Colors)",
			args: args{
				answer: Line{White, White, Green, Green},
				sub:    Line{Pink, Red, Blue, Yellow},
			},
			want: HitBlow{nhit: 0, nblow: 0},
		},
		{
			name: "1Hits/2Blows (2 Colors)",
			args: args{
				answer: Line{Red, Red, Yellow, Yellow},
				sub:    Line{Red, Yellow, Blue, Red},
			},
			want: HitBlow{nhit: 1, nblow: 2},
		},
		{
			name: "4Hits/0Blows (1 Color)",
			args: args{
				answer: Line{Pink, Pink, Pink, Pink},
				sub:    Line{Pink, Pink, Pink, Pink},
			},
			want: HitBlow{nhit: 4, nblow: 0},
		},
		{
			name: "1Hits/0Blows (1 Color)",
			args: args{
				answer: Line{Pink, Pink, Pink, Pink},
				sub:    Line{Yellow, Pink, Red, Blue},
			},
			want: HitBlow{nhit: 1, nblow: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcHB(&tt.args.answer, &tt.args.sub); !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("calcHB() = %v, want %v", *got, tt.want)
			}
		})
	}
}

func TestHitBlow_String(t *testing.T) {
	type fields struct {
		nhit  int
		nblow int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Test case 1",
			fields: fields{nhit: 3, nblow: 1},
			want:   "{3 hits, 1 blows}",
		},
		{
			name:   "Test case 2",
			fields: fields{nhit: 0, nblow: 4},
			want:   "{0 hits, 4 blows}",
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hb := HitBlow{
				nhit:  tt.fields.nhit,
				nblow: tt.fields.nblow,
			}
			if got := hb.String(); got != tt.want {
				t.Errorf("HitBlow.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_ShortString(t *testing.T) {
	tests := []struct {
		name  string
		color Color
		want  string
	}{
		{
			name:  "Test case 1",
			color: Blue,
			want:  "B",
		},
		{
			name:  "Test case 2",
			color: Red,
			want:  "R",
		},
		// Add more test cases as needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.ShortString(); got != tt.want {
				t.Errorf("Color.ShortString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_String(t *testing.T) {
	tests := []struct {
		name string
		line Line
		want string
	}{
		{
			name: "Test case 1",
			line: Line{Blue, Red, Green, Yellow},
			want: "BRGY",
		},
		{
			name: "Test case 2",
			line: Line{Pink, White, Blue, Red},
			want: "PWBR",
		},
		{
			name: "Test case 3",
			line: Line{White, White, White, White},
			want: "WWWW",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.String(); got != tt.want {
				t.Errorf("Line.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLine_GetLineId(t *testing.T) {
	tests := []struct {
		name string
		line Line
		want LineId
	}{
		{
			line: Line{Blue, Blue, Blue, Blue},
			want: LineId(0),
		},
		{
			line: Line{Green, Blue, White, Red},
			want: LineId((int(Green) * 6 * 6 * 6) + (int(Blue) * 6 * 6) + (int(White) * 6) + int(Red)),
		},
		{
			line: Line{White, White, White, White},
			want: LineId(5 + (5 * 6) + (5 * 6 * 6) + (5 * 6 * 6 * 6)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.GetLineId(); got != tt.want {
				t.Errorf("Line.GetLineId() = %v, want %v", got, tt.want)
			}
		})
	}
}
