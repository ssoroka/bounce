package levels

import (
	"encoding/json"
	"os"
)

func (g *Game) SaveState(filename string) error {
	b, err := json.Marshal(g)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

type typeChecker struct {
	Type string `json:"Type"`
}

func (g *Game) LoadState(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	g.Objects = nil // Clear existing objects
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(g); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (g *Game) UnmarshalJSON(data []byte) error {
	type Alias Game
	aux := struct {
		*Alias
		Objects []json.RawMessage `json:"Objects"`
	}{
		Alias: (*Alias)(g),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	for _, objData := range aux.Objects {
		var tc typeChecker
		if err := json.Unmarshal(objData, &tc); err != nil {
			return err
		}
		switch tc.Type {
		case "Circle":
			var c Circle
			if err := json.Unmarshal(objData, &c); err != nil {
				return err
			}
			g.Objects = append(g.Objects, &c)
		case "Boundary":
			var b Boundary
			if err := json.Unmarshal(objData, &b); err != nil {
				return err
			}
			g.Objects = append(g.Objects, &b)
		case "Cube":
			var c Cube
			if err := json.Unmarshal(objData, &c); err != nil {
				return err
			}
			g.Objects = append(g.Objects, &c)
		default:
			return nil // or return an error for unknown type
		}
	}
	return nil
}
