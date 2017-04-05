#version 330
uniform sampler2D tex;
uniform int textureUsed;
in vec2 fragTexCoord;
in vec4 fragColor;
out vec4 outputColor;

void main() {
    vec4 tmp;
    if (textureUsed == 1) {
        tmp = texture2D(tex, fragTexCoord);
    } else {
        tmp = vec4(1, 1, 1, 1);
    }
    outputColor = tmp * fragColor;
}
