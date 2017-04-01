#version 330

uniform mat2 normalMatrix;

in vec2 vert;
in vec2 vertTexCoord;
in float rotation;
in vec2 centerPoint;
in vec4 inColor;

out vec2 fragTexCoord;
out vec4 fragColor;

// standard openGL output
//out gl_PerVertex {
//    vec4  gl_Position;
//    float gl_PointSize;
//    float gl_ClipDistance[];
//};

vec2 rotate2D(in vec2 center, in float rotation, in vec2 vert) {
    mat2 rotationMatrix = mat2(
        cos(rotation), sin(rotation), // first column!
        -sin(rotation), cos(rotation)
    );
    return center + rotationMatrix * (vert - center);
}

void main() {
    fragColor = inColor;
    fragTexCoord = vertTexCoord;
    vec2 rotated = rotate2D(normalMatrix*centerPoint, rotation, vert);
    gl_Position = vec4( normalMatrix * rotated, 0, 1);
    //gl_Position = vec4( vert.x, vert.y, 0, 1);
    //gl_Position  = vec4(0.5, 0, 0, 1);
    gl_PointSize=5;
}
