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
	var vts []Tex

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
		case strings.HasPrefix(line, "vt "):
			var vt Tex
			matches, err := fmt.Fscanf(strings.NewReader(line), "vt %f %f", &vt.u, &vt.v)
			if err != nil || matches != 2 {
				log.Fatalf("vertex: only %d matches on line %q\n", matches, line)
			}
			vt.v = 1 - vt.v
			vts = append(vts, vt)
		case strings.HasPrefix(line, "f "):
			var vertexIndices [3]int
			var normalIndices [3]int
			var textureIndices [3]int
			f := strings.NewReader(line)

			if !strings.Contains(line, "/") {
				matches, err := fmt.Fscanf(f, "f %d %d %d", &vertexIndices[0], &vertexIndices[1], &vertexIndices[2])
				if err != nil || matches != 3 {
					log.Fatalf("face: only %d matches on line %q\n", matches, line)
				}

				// Append Face
			} else {
				matches, err := fmt.Fscanf(f, "f %d/%d/%d %d/%d/%d %d/%d/%d",
					&vertexIndices[0], &textureIndices[0], &normalIndices[0],
					&vertexIndices[1], &textureIndices[1], &normalIndices[1],
					&vertexIndices[2], &textureIndices[2], &normalIndices[2],
				)
				if err != nil || matches != 9 {
					log.Fatalf("face: only %d matches on line %q\n", matches, line)
				}
			}

			triangle := Triangle{
				color: White,
				texcoords: [3]Tex{
					vts[textureIndices[0]-1],
					vts[textureIndices[1]-1],
					vts[textureIndices[2]-1],
				},
				vertices: [3]Vec3{
					vertices[vertexIndices[0]-1],
					vertices[vertexIndices[1]-1],
					vertices[vertexIndices[2]-1],
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
