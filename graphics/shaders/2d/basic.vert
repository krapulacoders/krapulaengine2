#version 330
in vec2 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = vec4(vert, 1, 1);
}