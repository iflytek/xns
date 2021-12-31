package fastserver
/*
type DocParam struct {
	Name     string
	Type     string
	Required string
	Default  string
	Desc     string
}

type Document struct {
	Name    string
	Method  string
	Route   string
	Desc    string
	Params  []*DocParam
	Example string
    RespParams
}
type ApiGroup struct {
	apis      []*Api
	docParams []*DocParam
	Docs      []*Document
}
 */

const apiMarkDownTemplate = `
## api 文档
{{ range $api:=.Docs -}}
#### {{$api.Name}}
**{{$api.Desc}}**
{{$api.Method}} {{$api.Route}}
参数名称|参数类型|是否必传|默认值|描述
---|---|---|----|---
{{range $p:=$api.Params -}}
{{$p.Name}}|{{$p.Type}}|{{$p.Required}}|{{$p.Default}}|{{$p.Desc}}{{println}}
{{- end }}
**请求json示例**
....
{{$api.Example}}
····
{{- end }}

`
const apiHtmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.ServiceName}} api 文档</title>
</head>
<style type="text/css">
html {
    font-family: sans-serif;
    -ms-text-size-adjust: 100%;
    -webkit-text-size-adjust: 100%;
}
body {
    margin: 10px;
}
table {
    border-collapse: collapse;
    border-spacing: 4;
    border: 1px solid #cbcbcb;
  border-left: 1px solid #cbcbcb;
    border-width: 0 0 0 1px;
    font-size: inherit;
    margin: 0;
    overflow: visible;
    padding: .5em 1em;
    border-collapse: collapse;
    border-spacing: 0;
    empty-cells: show;
    border: 1px solid #cbcbcb;
}

td,th {
    padding: 5px;
}

textarea{
 padding: 0;
 height: auto;
 width: 400px;
}

input {
border-widht:3px;
font-size:18px;
width: 350px;
}


.header{border-bottom:1px solid #ccc;margin-bottom:5px;}
.MainContainer{display:flex;min-width:960px;max-width:1600px;}
.sidebar{width:280px;height:calc( 100% - 30px );margin-right:-180px;border-right:1px solid #ccc;min-height:500px;padding:5px; position:fixed;overflow:auto}
.main{margin-left:310px;padding:5px;}
.content{padding:0 10px;}

</style>
<body>


 <div class="page">
        <div class="MainContainer">
            <div class="sidebar">
            <h3 id="indexes"> 目录 </h3>
{{ range $i,$api:=.Docs -}}
<a href="#index_{{$i}}">{{$i}}__{{$api.Name}}</a><br>
{{ end }}
            </div>
            <div  id="main1" class="main">
                
{{ range $i,$api:=.Docs -}}
<h3 id="index_{{$i}}"> 1.{{$i}} {{$api.Name}} </h3>
<h4>{{$api.Desc}}</h4>
<h4>{{$api.Method}} &nbsp;<input id='input{{$i}}' value='{{$api.Route}}'></input></h4>
<h4>Content-Type: {{$api.ContentType}}</h4>
<h4>请求参数：</h4>
<table border="1">
<tr>
<th>name</th>
<th>type</th>
<th>required</th>
<th>default</th>
<th>constraints</th>
<th>from</th>
<th>description</th>
</tr>

{{range $p:=$api.Params -}}
<tr>
<td> {{$p.Name}} </td>
<td> {{$p.Type}} </td>
<td> {{$p.Required}} </td>
<td> {{$p.Default}} </td>
<td> {{$p.Constraints}} </td>
<td> {{$p.From}} </td>
<td> {{$p.Desc}} </td>
</tr>
{{- end }}
</table>

<h5>请求示例：</h5>
<textarea autoHeight cols="30" id="area{{$i}}">
{{$api.Example}}
</textarea>
<button onClick="request('{{$api.Method}}','{{$api.Route}}','{{$i}}','{{$api.ContentType}}')">发送请求</button>
<br>
<a href="#indexes">返回目录</a><br>
<a href="#status_anno">响应status解析：</a>

<h5>响应示例：</h5>
<textarea id="resp{{$i}}" autoHeight cols="30">{{$api.RespExample}}</textarea>

{{- end }}


<h3>响应status解析:</h3>
<table border="1" id="status_anno">
<tr>
<th>statusCode</th>
<th>desc</th>
</tr>
<tr>
<td>200</td>
<td>请求成功</td>
</tr>

<tr>
<td>404</td>
<td>没有找到目标</td>
</tr>
<tr>

<tr>
<td>400</td>
<td>请求参数有误，根据返回的消息,修正请求参数</td>
</tr>

<td>500</td>
<td>服务内部错误</td>
</tr>
</table>
            </div>           
        </div>
</div>





</body>
</html>
<script src="https://cdn.bootcss.com/jquery/3.4.1/jquery.js"></script>
<script>

$(function(){
    $.fn.autoHeight = function(){
        function autoHeight(elem){
            elem.style.height = 'auto';
            elem.scrollTop = 0; //防抖动
            elem.style.height = elem.scrollHeight+2 + 'px';
        }
        this.each(function(){
            autoHeight(this);
            $(this).on('input', function(){
                autoHeight(this);
            });
 			$(this).on('change', function(){
                autoHeight(this);
            });
        });
    }
    $('textarea[autoHeight]').autoHeight();
})

function request(method,route,name,contentType){
        $.ajax({
            url: $('#input'+name).val(),
            type:method,
            data:$('#area'+name).val(),
			contentType:contentType,
            success: function (result) {
                //console.log(result)
				t = $('#resp'+name)
				t.val(result)

            },
			error: function (result){
				t = $('#resp'+name)
				t.val(result.responseText)
				
			}
        });
}

</script>

`
