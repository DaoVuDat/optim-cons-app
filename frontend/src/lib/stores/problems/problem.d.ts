
export interface Coordinate {
    X: number
    Y: number
}

export interface Facility {
    Coordinate: Coordinate
    Rotation:   boolean
    Length:     number
    Width:      number
    IsFixed:    boolean
    Symbol:     string
    Name:       string
    IsLocatedAt: string
}