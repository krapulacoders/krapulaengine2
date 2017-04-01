#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
in vec4 fragColor;
out vec4 outputColor;

void main() {
    vec4 tmp = texture2D(tex, fragTexCoord);
    outputColor = fragColor;
}
