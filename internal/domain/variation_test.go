package domain

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariation_Linear(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		x        float64
		y        float64
		expected [2]float64
	}{
		{"zero", 0, 0, [2]float64{0, 0}},
		{"positive", 1, 2, [2]float64{1, 2}},
		{"negative", -1, -2, [2]float64{-1, -2}},
		{"mixed", 1.5, -2.5, [2]float64{1.5, -2.5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Linear(tt.x, tt.y)
			assert.InDelta(t, tt.expected[0], x, 1e-10, "x coordinate")
			assert.InDelta(t, tt.expected[1], y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Sinusoidal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"pi/2", math.Pi / 2, math.Pi / 2},
		{"pi", math.Pi, math.Pi},
		{"negative", -math.Pi / 2, -math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Sinusoidal(tt.x, tt.y)
			assert.InDelta(t, math.Sin(tt.x), x, 1e-10, "x coordinate")
			assert.InDelta(t, math.Sin(tt.y), y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Spherical(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit circle", 1, 0},
		{"diagonal", 1, 1},
		{"small values", 0.1, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Spherical(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r2 := tt.x*tt.x + tt.y*tt.y
				expectedX := tt.x / r2
				expectedY := tt.y / r2
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Swirl(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Swirl(tt.x, tt.y)
			r2 := tt.x*tt.x + tt.y*tt.y
			expectedX := tt.x*math.Sin(r2) - tt.y*math.Cos(r2)
			expectedY := tt.x*math.Cos(r2) + tt.y*math.Sin(r2)
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Horseshoe(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Horseshoe(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
				expectedX := 1 / r * (tt.x - tt.y) * (tt.x + tt.y)
				expectedY := 1 / r * 2 * tt.x * tt.y
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Disc(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Disc(tt.x, tt.y)
			r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
			ttt := math.Atan2(tt.y, tt.x)
			expectedX := ttt / math.Pi * math.Sin(math.Pi*r)
			expectedY := ttt / math.Pi * math.Cos(math.Pi*r)
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Diamond(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Diamond(tt.x, tt.y)
			r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
			ttt := math.Atan2(tt.y, tt.x)
			expectedX := math.Sin(ttt) * math.Cos(r)
			expectedY := math.Sin(r) * math.Cos(ttt)
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Ex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Ex(tt.x, tt.y)
			r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
			ttt := math.Atan2(tt.y, tt.x)
			p0 := math.Sin(ttt + r)
			p1 := math.Cos(ttt - r)
			expectedX := r * (p0*p0*p0 + p1*p1*p1)
			expectedY := r * (p0*p0*p0 - p1*p1*p1)
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Polar(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Polar(tt.x, tt.y)
			r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
			ttt := math.Atan2(tt.y, tt.x)
			expectedX := ttt / math.Pi
			expectedY := r - 1
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Handkerchief(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Handkerchief(tt.x, tt.y)
			r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
			ttt := math.Atan2(tt.y, tt.x)
			expectedX := r * math.Sin(ttt+r)
			expectedY := r * math.Cos(ttt-r)
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Heart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Heart(tt.x, tt.y)
			r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
			ttt := math.Atan2(tt.y, tt.x)
			expectedX := r * math.Sin(ttt*r)
			expectedY := -r * math.Cos(ttt*r)
			assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
			assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Spiral(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Spiral(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
				ttt := math.Atan2(tt.y, tt.x)
				invR := 1.0 / r
				expectedX := invR * (math.Cos(ttt) + math.Sin(r))
				expectedY := invR * (math.Sin(ttt) - math.Cos(r))
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Hyperbolic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Hyperbolic(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
				ttt := math.Atan2(tt.y, tt.x)
				expectedX := math.Sin(ttt) / r
				expectedY := r * math.Cos(ttt)
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Bent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		x        float64
		y        float64
		expected [2]float64
	}{
		{"both positive", 1, 2, [2]float64{1, 2}},
		{"x negative", -1, 2, [2]float64{2 * -1, 2}},
		{"y negative", 1, -2, [2]float64{1, -2 / 2}},
		{"both negative", -1, -2, [2]float64{2 * -1, -2 / 2}},
		{"zero", 0, 0, [2]float64{0, 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Bent(tt.x, tt.y)
			assert.InDelta(t, tt.expected[0], x, 1e-10, "x coordinate")
			assert.InDelta(t, tt.expected[1], y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Fisheye(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Fisheye(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
				re := 2.0 / (r + 1)
				expectedX := re * tt.y
				expectedY := re * tt.x
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Eyefish(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Eyefish(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
				re := 2.0 / (r + 1)
				expectedX := re * tt.x
				expectedY := re * tt.y
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Bubble(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Bubble(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r2 := tt.x*tt.x + tt.y*tt.y
				re := 4.0 / (r2 + 4)
				expectedX := re * tt.x
				expectedY := re * tt.y
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Cylinder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"pi/2", math.Pi / 2, 1},
		{"pi", math.Pi, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Cylinder(tt.x, tt.y)
			assert.InDelta(t, math.Sin(tt.x), x, 1e-10, "x coordinate")
			assert.InDelta(t, tt.y, y, 1e-10, "y coordinate")
		})
	}
}

func TestVariation_Tangent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"small", 0.5, 0.5},
		{"pi/4", math.Pi / 4, math.Pi / 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Tangent(tt.x, tt.y)
			cosY := math.Cos(tt.y)
			if math.Abs(cosY) < eps {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				expectedX := math.Sin(tt.x) / cosY
				expectedY := math.Tan(tt.y)
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Cross(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 1},
		{"different", 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Cross(tt.x, tt.y)
			diff2 := tt.x*tt.x - tt.y*tt.y
			if math.Abs(diff2) < eps {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				s := math.Sqrt(1.0 / (diff2 * diff2))
				expectedX := s * tt.x
				expectedY := s * tt.y
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestVariation_Power(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		x    float64
		y    float64
	}{
		{"zero", 0, 0},
		{"unit", 1, 0},
		{"diagonal", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := Power(tt.x, tt.y)
			if tt.x == 0 && tt.y == 0 {
				assert.Equal(t, 0.0, x)
				assert.Equal(t, 0.0, y)
			} else {
				r := math.Sqrt(tt.x*tt.x + tt.y*tt.y)
				ttt := math.Atan2(tt.y, tt.x)
				rsth := math.Pow(r, math.Sin(ttt))
				expectedX := rsth * math.Cos(ttt)
				expectedY := rsth * math.Sin(ttt)
				assert.InDelta(t, expectedX, x, 1e-10, "x coordinate")
				assert.InDelta(t, expectedY, y, 1e-10, "y coordinate")
			}
		})
	}
}

func TestGetVariation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		variation string
		wantErr   bool
	}{
		{"linear", "linear", false},
		{"sinusoidal", "sinusoidal", false},
		{"spherical", "spherical", false},
		{"swirl", "swirl", false},
		{"horseshoe", "horseshoe", false},
		{"polar", "polar", false},
		{"handkerchief", "handkerchief", false},
		{"heart", "heart", false},
		{"disc", "disc", false},
		{"spiral", "spiral", false},
		{"hyperbolic", "hyperbolic", false},
		{"diamond", "diamond", false},
		{"ex", "ex", false},
		{"bent", "bent", false},
		{"fisheye", "fisheye", false},
		{"eyefish", "eyefish", false},
		{"bubble", "bubble", false},
		{"cylinder", "cylinder", false},
		{"tangent", "tangent", false},
		{"cross", "cross", false},
		{"power", "power", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v, err := GetVariation(tt.variation)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, v)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, v)
			}
		})
	}
}

func TestVariationList_Add(t *testing.T) {
	t.Parallel()
	t.Run("add valid variation", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		v, _ := GetVariation("linear")
		err := list.Add(v, 1.0)
		assert.NoError(t, err)
		assert.Equal(t, 1, list.Len())
		assert.InDelta(t, 1.0, list.TotalWeight(), 1e-10)
	})

	t.Run("add multiple variations", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		v1, _ := GetVariation("linear")
		v2, _ := GetVariation("sinusoidal")
		err1 := list.Add(v1, 1.0)
		err2 := list.Add(v2, 2.0)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 2, list.Len())
		assert.InDelta(t, 3.0, list.TotalWeight(), 1e-10)
	})

	t.Run("add negative weight", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		v, _ := GetVariation("linear")
		err := list.Add(v, -1.0)
		assert.Error(t, err)
		assert.Equal(t, ErrNegativeVariationWeight, err)
		assert.Equal(t, 0, list.Len())
	})
}

func TestVariationList_Apply(t *testing.T) {
	t.Parallel()
	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		x, y := list.Apply(1.0, 2.0)
		assert.Equal(t, 1.0, x)
		assert.Equal(t, 2.0, y)
	})

	t.Run("single variation", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		v, _ := GetVariation("linear")
		_ = list.Add(v, 1.0)
		x, y := list.Apply(1.0, 2.0)
		assert.Equal(t, 1.0, x)
		assert.Equal(t, 2.0, y)
	})

	t.Run("multiple variations with weights", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		v1, _ := GetVariation("linear")
		v2, _ := GetVariation("sinusoidal")
		_ = list.Add(v1, 1.0)
		_ = list.Add(v2, 1.0)
		x, y := list.Apply(1.0, 2.0)
		expectedX := (1.0 + math.Sin(1.0)) / 2.0
		expectedY := (2.0 + math.Sin(2.0)) / 2.0
		assert.InDelta(t, expectedX, x, 1e-10)
		assert.InDelta(t, expectedY, y, 1e-10)
	})

	t.Run("zero total weight", func(t *testing.T) {
		t.Parallel()
		list := NewVariationList()
		v, _ := GetVariation("linear")
		_ = list.Add(v, 0.0)
		x, y := list.Apply(1.0, 2.0)
		assert.Equal(t, 1.0, x)
		assert.Equal(t, 2.0, y)
	})
}

func TestVariationList_Variations(t *testing.T) {
	t.Parallel()
	list := NewVariationList()
	v1, _ := GetVariation("linear")
	v2, _ := GetVariation("sinusoidal")
	_ = list.Add(v1, 1.0)
	_ = list.Add(v2, 2.0)

	variations, weights := list.Variations()
	require.Equal(t, 2, len(variations))
	require.Equal(t, 2, len(weights))
	assert.Equal(t, 1.0, weights[0])
	assert.Equal(t, 2.0, weights[1])
}

func TestGetVariationNames(t *testing.T) {
	t.Parallel()
	names := GetVariationNames()
	assert.Greater(t, len(names), 0)

	for _, name := range names {
		v, err := GetVariation(name)
		assert.NoError(t, err, "variation %s should be valid", name)
		assert.NotNil(t, v, "variation %s should not be nil", name)
	}
}
