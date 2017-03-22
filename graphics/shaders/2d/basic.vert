#version 330

uniform mat2 normalMatrix;

in vec2 vert;
in vec2 vertTexCoord;
in float rotation;
in vec2 centerPoint;
in vec4 inColor;

out vec2 fragTexCoord;
out vec4 outColor;

// standard openGL output
//out gl_PerVertex {
//    vec4  gl_Position;
//    float gl_PointSize;
//    float gl_ClipDistance[];
//};

vec2 rotate2D(vec2 center, float rotation, vec2 vert) {
    mat2 rotationMatrix = mat2(
        cos(rotation), sin(rotation), // first column!
        -sin(rotation), cos(rotation)
    );
    return center + rotationMatrix * (vert - center);
}

void main() {
    outColor = inColor;
    fragTexCoord = vertTexCoord;
    vec2 rotated = rotate2D(centerPoint, rotation, vert);
    gl_Position = vec4( normalMatrix * rotated, 1, 1);
}
