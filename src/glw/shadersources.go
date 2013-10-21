package glw

var SHADER_SOURCES = map[ShaderRef]ShaderSeed{
	VSH_POS3: ShaderSeed{Type: VERTEX_SHADER, Source: `
#version 330 core
// Transmits the position to the fragment shader.

layout(std140) uniform GlobalMatrices
{
    mat4 projection_matrix;
    mat4 view_matrix;
};
uniform mat4 model_matrix;
in vec3 vpos;
out vec3 fpos;
void main(){
	vec4 temp = model_matrix * vec4(vpos, 1.0);
	temp = view_matrix * temp;
	gl_Position = projection_matrix * temp;
	fpos = gl_Position.xyz;
}`},
	VSH_COL3: ShaderSeed{Type: VERTEX_SHADER, Source: `
#version 330 core
// Transmits the color to the fragment shader.

layout(std140) uniform GlobalMatrices
{
    mat4 projection_matrix;
    mat4 view_matrix;
};
uniform mat4 model_matrix;
in vec3 vpos;
in vec3 vcol;
out vec3 fcol;
void main(){
	vec4 temp = model_matrix * vec4(vpos, 1.0);
	temp = view_matrix * temp;
	gl_Position = projection_matrix * temp;
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
