#version 330

uniform mat4 normalMatrix;
uniform mat4 modelMatrix;

in vec3 vert;
in vec2 vertTexCoord;
in float rotation;
in vec3 centerPoint;
in vec4 inColor;

out vec2 fragTexCoord;
out vec4 fragColor;

// standard openGL output
//out gl_PerVertex {
//    vec4  gl_Position;
//    float gl_PointSize;
//    float gl_ClipDistance[];
//};

vec3 rotateAroundZ(in vec3 center, in float rotation, in vec3 vert) {
    mat3 rotationMatrix = mat3(
        cos(rotation), sin(rotation), 0, // first column!
        -sin(rotation), cos(rotation), 0,
        0, 0, 1
    );
    return center + rotationMatrix * (vert - center);
}

vec3 rotateAroundX(in vec3 center, in float rotation, in vec3 vert) {
    mat3 rotationMatrix = mat3(
        1, 0, 0,
        0, cos(rotation), sin(rotation), // first column!
        0, -sin(rotation), cos(rotation)
    );
    return center + rotationMatrix * (vert - center);
}

vec3 rotateAroundY(in vec3 center, in float rotation, in vec3 vert) {
    mat3 rotationMatrix = mat3(
        cos(rotation), 0, -sin(rotation), // first column!
        0, 1, 0,
        sin(rotation), 0, cos(rotation)
    );
    return center + rotationMatrix * (vert - center);
}

void main() {
    fragColor = inColor;
    fragTexCoord = vertTexCoord;
    vec3 rotated = rotateAroundZ(centerPoint, rotation, vert);
    gl_Position = normalMatrix * modelMatrix * vec4(rotated, 1);
    gl_PointSize= 10 * ( 1.1 - gl_Position.z);
}
