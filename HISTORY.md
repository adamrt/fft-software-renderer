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

## Design Decisions

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

## Daily Notes
### 2022-12-12

Today I ready through chapter 1 and part of 2 of [3D Math Primer for
Graphics and Game
Development](https://gamemath.com/book/index.html). I've done lots of
3D studying over the past year, but I've skimmed this (which I found
yesterday) and the way its written really clicks with me. I'm hoping
to get a better understanding as I build out this engine.


**Result**: Draw a cube points on the screen with orthographic
projection and rotate them.


#### Mathmatical notes

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



### 2022-12-11

- Add Window struct for wrapping the SDL_Window and SDL_Renderer.
- Add Engine struct for holding a window and standard engine loop.

  A typical engine loop has `processInput()`, `update()` and
  `render()` methods.

- Add Color struct as 4 bytes representing RGBA.

  NOTE: The SDL pixel format we used is ABGR8888 because it reads the bytes
  in big endian instead of little endian.

- Add forground and background SDL textures.

- Add fullscreen and windowed modes

**Result**: Opening a window and fill with background and draw a
pixel. Escape to quit.
