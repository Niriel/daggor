package glw

var SHADER_SOURCES = map[ShaderRef]ShaderSeed{
	VSH_POS3: ShaderSeed{Type: VERTEX_SHADER, Source: `
#version 330 core
// Transmits the position to the fragment shader.
in vec3 vpos;
out vec3 fpos;
uniform mat4 mvp;
void main(){
	gl_Position = mvp * vec4(vpos, 1.0);
	fpos = gl_Position.xyz;
}`},
	VSH_COL3: ShaderSeed{Type: VERTEX_SHADER, Source: `
#version 330 core
// Transmits the color to the fragment shader.
in vec3 vpos;
in vec3 vcol;
uniform mat4 mvp;
out vec3 fcol;
void main(){
	gl_Position = mvp * vec4(vpos, 1.0);
	fcol = vcol;
}`},
	FSH_ZRED: ShaderSeed{Type: FRAGMENT_SHADER, Source: `
#version 330 core
// Red fragment, the farther the darker (from z coordinate).
in vec3 fpos;
out vec3 color;

void main(){
	color = vec3(1.0 - fpos.z *.1, 0, 0);
}`},
	FSH_ZGREEN: ShaderSeed{Type: FRAGMENT_SHADER, Source: `
#version 330 core
// Green fragment, the farther the darker (from z coordinate).
in vec3 fpos;
out vec3 color;

void main(){
	color = vec3(0, 1.0 - fpos.z *.1, 0);
}`},
	FSH_VCOL: ShaderSeed{Type: FRAGMENT_SHADER, Source: `
#version 330 core
// Takes color from vertex.
in vec3 fcol;
out vec3 color;

void main(){
	color = fcol;
}`},
}
