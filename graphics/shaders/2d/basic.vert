#version 330

uniform Matrices {
    mat2 normalMatrix;
}

in vec2 vert;
in vec2 vertTexCoord;
in float rotation;
in vec2 centerPoint;
in vec4 inColor;

out vec2 fragTexCoord;
out vec4 outColor;

// standard openGL output
out gl_PerVertex {
    vec4  gl_Position;
    float gl_PointSize;
    float gl_ClipDistance[];
};

void main() {
    outColor = inColor;
    fragTexCoord = vertTexCoord;
    gl_Position = vec4( vec2(normalMatrix * rotate2D(centerPoint, rotation, vert), 1, 1);
}

vec2 rotate2D(vec2 center, float rotation, vec2 vert) {
    mat2 rotationMatrix = mat2(
        cos(rotation), sin(rotation) // first column!
        -sin(rotation), cos(rotation)
    )
    return center + rotationMatrix * (vert - center)
}