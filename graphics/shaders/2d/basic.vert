#version 330

uniform mat4 normalMatrix;
uniform mat4 modelMatrix;
uniform ivec2 rotationMode;

in vec3 vert;
in vec2 vertTexCoord;
in vec2 angles;
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

vec3 rotateAroundZ(in vec3 center, in float angle, in vec3 vert) {
    mat3 rotationMatrix = mat3(
        cos(angle), sin(angle), 0, // first column!
        -sin(angle), cos(angle), 0,
        0, 0, 1
    );
    return center + rotationMatrix * (vert - center);
}

vec3 rotateAroundX(in vec3 center, in float angle, in vec3 vert) {
    mat3 rotationMatrix = mat3(
        1, 0, 0,
        0, cos(angle), sin(angle), // first column!
        0, -sin(angle), cos(angle)
    );
    return center + rotationMatrix * (vert - center);
}

vec3 rotateAroundY(in vec3 center, in float angle, in vec3 vert) {
    mat3 rotationMatrix = mat3(
        cos(angle), 0, -sin(angle), // first column!
        0, 1, 0,
        sin(angle), 0, cos(angle)
    );
    return center + rotationMatrix * (vert - center);
}

vec3 rotate(in vec3 center, in float angle, in vec3 vert, int mode) {
    if (mode == 1) {
        return rotateAroundX(center, angle, vert);
    } else if (mode == 2) {
        return rotateAroundY(center, angle, vert);
    } else if (mode == 3) {
        return rotateAroundZ(center, angle, vert);
    } else {
        return vert;
    }
}

void main() {
    fragColor = inColor;
    fragTexCoord = vertTexCoord;
    vec3 rotated = rotate(centerPoint, angles.x, vert, rotationMode.x);
    rotated = rotate(centerPoint, angles.y, rotated, rotationMode.y);
    gl_Position = normalMatrix * modelMatrix * vec4(rotated, 1);
    gl_PointSize= 10 * ( 1.1 - gl_Position.z);
}
