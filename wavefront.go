// This file is for loading wavefront obj files as meshes.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func NewMeshFromObj(objFilename string) Mesh {
	objFile, err := os.Open(objFilename)
	if err != nil {
		panic(err)
	}
	defer objFile.Close()

	mesh := NewMesh()

	vertices := []Vec3{}
	scanner := bufio.NewScanner(objFile)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			var v Vec3
			matches, err := fmt.Fscanf(strings.NewReader(line), "v %f %f %f", &v.x, &v.y, &v.z)
			if err != nil || matches != 3 {
				log.Fatalf("vertex: only %d matches on line %q\n", matches, line)
			}
			vertices = append(vertices, v)
		case strings.HasPrefix(line, "f "):
			var vertex_indices [3]int
			var normal_indices [3]int
			var texture_indices [3]int
			f := strings.NewReader(line)

			if !strings.Contains(line, "/") {
				matches, err := fmt.Fscanf(f, "f %d %d %d", &vertex_indices[0], &vertex_indices[1], &vertex_indices[2])
				if err != nil || matches != 3 {
					log.Fatalf("face: only %d matches on line %q\n", matches, line)
				}

				// Append Face
			} else {
				matches, err := fmt.Fscanf(f, "f %d/%d/%d %d/%d/%d %d/%d/%d",
					&vertex_indices[0], &texture_indices[0], &normal_indices[0],
					&vertex_indices[1], &texture_indices[1], &normal_indices[1],
					&vertex_indices[2], &texture_indices[2], &normal_indices[2],
				)
				if err != nil || matches != 9 {
					log.Fatalf("face: only %d matches on line %q\n", matches, line)
				}
			}

			triangle := Triangle{
				color: White,
				vertices: [3]Vec3{
					vertices[vertex_indices[0]-1],
					vertices[vertex_indices[1]-1],
					vertices[vertex_indices[2]-1],
				},
			}
			mesh.triangles = append(mesh.triangles, triangle)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return mesh
}
