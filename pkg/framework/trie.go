package framework

import (
    "strings"
)

// node repräsentiert ein Segment der URL (z.B. "users" oder ":id")
type node struct {
    pattern  string           // Das Segment, z.B. "users" oder ":id"
    children []*node          // Kind-Knoten
    isWild   bool             // True, wenn es ein Parameter ist (beginnt mit :)
    handler  any		      // Der Handler, falls dies das Ende einer Route ist
}

// Helper: Fügt ein Kind hinzu, falls es noch nicht existiert
func (n *node) insertChild(pattern string, isWild bool) *node {
    for _, child := range n.children {
        if child.pattern == pattern {
            return child
        }
    }
    child := &node{pattern: pattern, isWild: isWild}
    n.children = append(n.children, child)
    return child
}

// insert traversiert den Baum und erstellt neue Nodes wo nötig.
func (n *node) insert(method, path string, handler any) {
    // 1. Pfad bereinigen und in Teile zerlegen
    // "/api/users/:id" -> ["api", "users", ":id"]
    parts := parsePath(path)

    current := n
    for _, part := range parts {
        // Prüfen ob Wildcard
        isWild := len(part) > 0 && part[0] == ':' || part[0] == '*'
        
        // Kind finden oder erstellen
        // Hinweis: Wir ignorieren hier Konflikte (z.B. static vs wild auf gleicher Ebene) für die Einfachheit
        current = current.insertChild(part, isWild)
    }

    // Am Ziel-Knoten den Handler speichern
    current.handler = handler
}

// parsePath zerlegt "/a/b/c" in ["a", "b", "c"]
func parsePath(path string) []string {
    vs := strings.Split(path, "/")
    parts := make([]string, 0)
    for _, item := range vs {
        if item != "" {
            parts = append(parts, item)
        }
    }
    return parts
}

// search sucht den Node für einen gegebenen Pfad und liefert Parameter zurück.
func (n *node) search(parts []string) (*node, map[string]string) {
    params := make(map[string]string)
    current := n

    for _, part := range parts {
        var next *node

        // Wir suchen in den Kindern nach einem Match
        for _, child := range current.children {
            if child.pattern == part || child.isWild {
                // Wenn Wildcard, speichern wir den Parameter!
                if child.isWild {
                    // child.pattern ist z.B. ":id", wir schneiden ":" ab
                    params[child.pattern[1:]] = part
                }
                next = child
                break
            }
        }

        if next == nil {
            return nil, nil // Sackgasse, Route nicht gefunden
        }
        current = next
    }

    // Haben wir einen Handler am Ende?
    if current.handler == nil {
        return nil, nil
    }

    return current, params
}