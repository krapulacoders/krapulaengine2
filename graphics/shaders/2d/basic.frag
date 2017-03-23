#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
in vec4 outColor;
out vec4 outputColor;
void main() {
    //outputColor = texture2D(tex, fragTexCoord);
    outputColor = outColor;
}
