# Heretic Project History

This is a log of the development and learning process for Heretic.

Once I realized you could extract level data from the original Final
Fantasy Tactics PS1 ISO I knew I had to build a simple map viewer. At
the time I was already a developer for many years but didn't have a
single clue about games or 3D development. I don't have a desire to
develop games to be creative, but I'm very interested in the inner
workings of 3D rendering.

This was in August 21, 2021 and I started the process of learning C
then C++. I took multiple courses and learned a ton along the way. I
built many small 3D renderers in C, C++, Go and Rust. Additionally I
built libraries in all of those languages to extract data from the
Final Fantasy Tactics PS1 ISO.

On December 12, 2022 I decided it was time to start from scratch and
build the engine I wanted. After switching between C, C++ and Go many
times, unable to decide on what to use, I landed on Go. I think Go is
the easiest to read and write and the performance is acceptable. The
only negative is being able to run the application on a PS Vita or
other mobile device. But that is a secondary concern. It could always
be ported to a more portable language in the future.

The only dependency is SDL2. All rendering is software rasterized
without OpenGL.

This is the journal of development.

### Resources

I used a million sources for learning 3D but I think these were some
of the ones I referred to the most.

- [Pikuma 3D Computer Graphics Programming Course](https://pikuma.com/courses/learn-3d-computer-graphics-programming)
- [3D Math Primer for Graphics and Game Development.](https://gamemath.com/)
- [Foundations of Game Engine Development: Volume 1](https://foundationsofgameenginedev.com/)
- [Javidx9 3D Graphics Engine Series](https://www.youtube.com/watch?v=ih20l3pJoeU)

### Design Decisions

- **Left-handed coordinates system**

  We used a left-handed coordinate system since it feels more
  intuitive. Left-handed sytem mean the positive z-axis is forward.

- **Clockwise vertex winding**

  Due to the left-handed coordinates system we use clockwise vertex winding.
  We would order the vertices of this triangle `{1, 3, 2}`.

  ```
     1
    / \
   2___3
  ```

- **Row-major matrices**

  Using row-major is much easier for my brain to develop in since we
  can use.

  ```
  type Matrix [4][4]float64
  Matrix{
    {x,y,z,w} // Row 1
    {x,y,z,w} // Row 2
    {x,y,z,w} // Row 3
    {x,y,z,w} // Row 4
  }
  ```

- **Floating point precision**

  Currently we are using `float64` for everything. `float32` would be
  perfectly fine and use less memory, but Go's math lib uses `float64`
  for everything. Additionally, we use so few vertices for a FFT map
  that the overhead is that much. This could be changed in the future
  is memory usage is high.

# Daily Notes
## 2022-01-27

There was a memory leak due to the sdl_ttf textures and surfaces that
are created every frame. I now destroy the textures and free the
surfaces each frame.

In the future if would make much more sense to just create these and
only recreate/destroy/free them when the options change.

## 2022-12-30

Add directional and ambient lighting. I'm not sure if the coloring is
working accurately at this time though.

Instead of modifying vertices during map load to center/normalize
them, we just scale the map and change the camera to point to the
center.

## 2022-12-29

Change to reading bin files instead of iso. It seems like iso was just
happening to work even though it kinda shouldn't have. Instead of
starting at the sector and reading as much as we needed, we now read
and concatenate consecutive sector data (discarding sector headers).

## 2022-12-28

Fix gaps in textures due to rounder errors when drawing textured
triangles.

## 2022-12-24

Add dynamic menu with options and keymappings.

## 2022-12-19

Move the view and projection matrix into the camera and implement
proper zooming. Previously zoom was acheived by scaling the
model. That looked okay, but when you have multiple models it is
obvious that you are scaling. Now we have a zoom factor.

Add level backgrounds

## 2022-12-17

Read FFT ISO. Add texture mapping and orbital controls.

## 2022-12-16

Today I started rasterizing triangles (filling them). This uses the
flat-top, flat-bottom technique that splits the triangle in two (if
necessary) so there can be two separate triangles. One with a flat-top
and a flat-bottom.

I basically copied Gustavo's code for this. I find this bit
uninteresting for now and understand how it works. There are faster
algorithms, but I don't care for now. I'll be using textures for most
faces anyway.

Painters algorithm was used to help with triangle depth ordering. Its
a naive approach for depth sorting, but we don't have (and may not
need) a depth buffer.

I added simple lighting as well. I thought it wasn't working because
my dot products were so close to 1.0 (.9999123) that the color wasn't
changing. This was because I was using the dot product on the vertices
after projection instead of before. The dp needs to be calculated on
the _transformed_ vertices, before they are projected.

I found a bug in the MatrixRotationZ code. It was just a typo when I
set it up. The z rotation was noticably off once I was looking at a
cube.

There was also a bug in the backface culling code. I had originally
switched the order of the forumla thinking that I was compensating for
CW/CCW winding orders, but it turned out to already be correct the way
it was in the Stack Overflow answer.


### Results

- Backface culling
- Triangle rasterization (filling triangles)
- Painters algorithm
- Simple lighting


## 2022-12-15

Implemented the new backface culling technique. Its so simple that I
really like it. You do have to do the calculation _after_ projection
so you don't get to short circuit the projection math, but I think
that is fine since the actual calculation is lighter than the previous
dot product method.

## 2022-12-14

Okay, after thinking about yesterdays problem for a while, I realized
that the option (#3) of inverting the Y-axis is the correct
solution. This is what Gustavo suggested as his solution [1] as
well. When I first did this course about a year ago, I didn't like
this solution since I thought it was to compensate for wavefront
format (.obj files). I didn't want to change our engine for the
format, and instead wanted to change the format during
read/import.

Now with a little more understanding I realize that the wavefront file
is correct (`+Y` axis up) and our medium is whats incorrect (`+Y` axis
down). Our medium being the color buffer. Yesterday, I was thinking
this fix would throw off calculations (dot product, lighting, culling
etc), down the road, but I was wrong. It will all work as expected.

This solution is also better than the other 2 options I listed as
well. Both of those options essentially invert the color buffer after
drawing. We have two methods of "drawing": setting individual pixels
starting at top-left (0,0), and drawing lines/triangles that have 3D
coordinates. The 3D coordinates go through transformation, projection,
etc, then, with option #3, have their Y-axis inverted. This means our
3D stuff gets drawn correctly, but our pixel drawings still get to use
the logical top-left (0,0) 2D coordinates, instead of starting at (0,
HEIGHT-1) if we used the "invert-the-whole-image" method.

**Follow-up 1**: I wrote the above comment to Gustavo thinking it was
solved. But I realized that due to the left-handed rule, the Y-axis
might be off. When pointing your thumb in the direction of a positive
axis (ie right for `+X`), if you curl your fingers that is the
direction the object should rotate with _positive_ rotation. This
works with the above fix for the `X` and `Z` axis, but the `Y` axis is
reversed.

**Follow-up 2**: Turns out that the Y-axis was off due to an actual
bug in the `rotate_y` function. The `-sin(angle)` was in the wrong
place. I was looking at this function because I was thinking I could
just mess with it to get the results I wanted, only to realize it was
actually wrong. I knew it was wrong because of a sentence I remembered
from [Foundations of Game Engine Development: Volume
1](https://foundationsofgameenginedev.com/) on page 61. It explains
where the `-sin` should go. _"The negated sine function always appears
one row below and one column to the left, with wraparound, of the
entry containing the one in the matrix."_ Thanks for that!

### Backface culling

- [x] Fix backface culling in orthographic projection.

    Right now backface culling is working perfect in perspective
    projection. But during orthographic, the normal or the camera ray
    vector is incorrect or something similar. I'm not sure if there
    are separate ways to handle ortho/perspective culling.

    I considered ignoring this since we have so few triangles that
    backface culling isn't that much of a performance improvement. But
    then without understanding why its happening, whats the point of
    continuing, more things would crop up down the line. Plus I think
    in the original go engine I did, switching to orthographic
    actually had this problem, now that I think about it.

    **Update**: Okay so I found an entirely new way to do backface
    culling in addition to why our current way wasn't working.

    [Stack Overflow Question](https://gamedev.stackexchange.com/questions/203694/is-backface-culling-affected-by-differently-between-orthographic-and-perspective)

    1. My orthographic issue was because I was pointing from the
       camera position instead of the camera direction. Since
       orthographic has no visible angle the same was as perspective,
       the camera position wouldn't angle correctly.

    1. The above doesn't really matter because I now have a better way
       to do culling by checking the triangle winding order _after
       projection_. It so much cheaper to compute and simplifies the
       code massivly.

More to say here but gotta bail.





### Notes
- Vector subtractions returns a vector pointing to the first element.

  `v := camera.position.Sub(someVector)`. v will point from someVetor to camera.

- The reason we need Vec4 and Mat4x4 instead of Vec3 and Mat3x3 is for two reasons.

  1. Translation is not a linear transformation

     Translation changes the origin, unlike scale and rotation (which
     are linear transformations).

  1. For perspective projection we can store the original z value for
     perspective in the w component by placing `[3][2] = 1.0`.

  I think we could use 3x3 if we only did scale and rotation.


[1] Non-Free Link: [Inverted Vertical Screen Values](https://courses.pikuma.com/courses/take/learn-computer-graphics-programming/lessons/12296518-inverted-vertical-screen-values/discussions/5719619)


### Result
- Scale, Rotation and Translation matrices (World Matrix)
- Perspective projection
- Backface culling in perspective projection

  Orthographic culling isn't working correctly, but I think that
  matrices might solve this. **Update**: It didn't.

## 2022-12-13

Last night and this morning I've been struggling with some fundamental
issue that I think I've run into in all of my previous prototypes.

With a left-handed coordinate system, the `+X`, `+Y` and `+Z` go
right, up and away from the viewer, respectively.

In the 2D colorbuffer that we use to render, the `+X` and `+Y` go
right and *down*. For the color buffer this makes sense because the
colorbuffer starts at `0,0` on the top-left corner.

The problem is that when projecting onto the screen everything is
upside down. In the past I've always just inverted the Y-axis of the
mesh during import.

But now I'm realizing that during translation, this also causes
problem. If I increase the `Y` coordinate by +1 on any mesh, the mesh
will move down on the screen instead of up.

Also if I draw axes from the origin and pointing to their positive direction like so:

```
{0,0,0} -> {1,0,0} (X axis)
{0,0,0} -> {0,1,0} (Y axis)
{0,0,0} -> {0,0,1} (Z axis)
```

The `+Y` axis points down.


This can be visually "fixed" by doing one of the following.

1. Change `buffer[(width*y)+x] = color` to `buffer[(width*(height-y-1))+x] = color`

    This will draw the buffer upside down.

1. Use `SDL_RendererFlip` with `SDL_RenderCopyEx()` to flip the image.

   This way you will draw the same way, but use SDL to flip the image
before rendering it.

1. Invert the Y-axis of vertices when importing a mesh.

I'm not sure that these fix the problem or just mask it. Also I'm not
sure that once we start using translation matrices if the problems
will be fixed, I don't think so but I might beable to understand it more.

- Removed the vertex indicies from the Mesh and just use the vertices
  directly (even if duplicated). These maps are very small so
  duplicating vertices isn't a big memory loss for now, and it makes
  the code easier to follow.

### Result

None just researching

## 2022-12-12

Today I ready through chapter 1 and part of 2 of [3D Math Primer for
Graphics and Game
Development](https://gamemath.com/book/index.html). I've done lots of
3D studying over the past year, but I've skimmed this (which I found
yesterday) and the way its written really clicks with me. I'm hoping
to get a better understanding as I build out this engine.

- Draw a cube points on the screen
- Rotate points with rotation functions per axes
- Add variable FPS timestep
- Add smallest possible orthographic projects (disgard the z component).
- Add Mesh, Face and Triangle

  Faces are for indexing vertices instead of duplicating overlapping vertex data.
  Triangles are the collection of Vec2 points after vertex projection

- Render vertex dots of triangles
- Render lines of triangles (wireframe)

### Notes

My basic math skills are about 20 years old, so I'm taking notes I can
refer back to.


- **Types of numbers**:
  - Natural: Positive whole numbers (1,2,3,100). Sometimes includes zero.
  - Integers: Natural numbers and their negative counterparts (...,-2,-1,0,1,2,...)
  - Rational: Every integer and any fractional number (-3/7). Denominator can't be zero.
  - Real: Rational numbers that require infinite decimal places (Pi)
  - The study of natural numbers and integers is called discrete mathematics, and the study of real numbers is called continuous mathematics.

- **Summation notation**: `Σ` is Capital Greek Sima

Basically a for loop. The number on top (6 in this case) is the last
usage, not repititions. if i=0, it would do 7 loops, `a0 + ... + a6`.
if `i=1`, it would do 6 loops, `a1 + ... + a6`.


```
  6
  Σ  a𝒾 = a1 + a2 + a3 + a4 + a5 + a6
 𝒾=i
```

- **Product notation**: `𝚷` is Capital Pi
```
  n
  𝚷  a𝒾 = a1 * a2 * ... * an-1 * an
 𝒾=i
```

- **Interval Notation**:

```
[a,b] == a ≤ x ≤ b
(a,b) == a < x < b

```

- **Angles, Degrees, Radians**:

    Extra notes: [Intuitive Guide to Angles, Degrees and Radians](https://betterexplained.com/articles/intuitive-guide-to-angles-degrees-and-radians/)

  - Theta (ϑ) typically represents angles.
  - Angles are typically measured in `degrees` and `radians`.
  - Humans usually use degress. 1 degree is 1/360.
  - Mathematitcians usually use radians.
  - Sin function is defined in terms of radians
  - (is this right?) In a Left-handed coordinate system, positive rotation rotates clockwise when viewed from the positive end of the axis.
  - Degrees are from observers view point (center of circle?). Radians are from the movers viewpoint.

    *"Degrees measure angles by how far we tilted our heads. Radians measure angles by distance traveled."*

    But absolute distance isn't that useful, since going 10 miles is a different number of laps depending on the track.

    So we divide by radius to get a normalized angle
    ```
              distance traveled                       s
    radian =  -----------------    commonly->    ϑ = ---
                    radius                            r

    ```
    or `angle in radians (theta) is arc length (s) divided by radius (r)`

  - Sin

    - `x` is how far you traveled along a circle
    - `sin(x)` is how high on the circle you are


### Result

- Rotating 3D wireframe cube


## 2022-12-11

- Add Window struct for wrapping the SDL_Window and SDL_Renderer.
- Add Engine struct for holding a window and standard engine loop.

  A typical engine loop has `processInput()`, `update()` and
  `render()` methods.

- Add Color struct as 4 bytes representing RGBA.

  NOTE: The SDL pixel format we used is ABGR8888 because it reads the bytes
  in big endian instead of little endian.

- Add forground and background SDL textures.

- Add fullscreen and windowed modes

### Result

Opening a window and fill with background and draw a pixel. Escape to
quit.
