## Heretic Project History

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

### 2022-12-11

- Add Window struct for wrapping the SDL_Window and SDL_Renderer.
- Add Engine struct for holding a window and standard engine loop.

  A typical engine loop has `processInput()`, `update()` and
  `render()` methods.

- Add Color struct as 4 bytes representing RGBA.

  NOTE: The SDL pixel format we used is ABGR8888 because it reads the bytes
  in big endian instead of little endian.

Result: Opening a window and fill with any color. Escape to quit.
