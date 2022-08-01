var dialogOpen, form, table;
layui.use(['form'], function () {
    form = layui.form;
    form.on('submit(btn_add)', function (data) {
        ajaxAdd(data.field);
        return false;
    });
    form.on('submit(btn_update)', function (data) {
        ajaxUpdate(data.field);
        return false;
    });
});
layui.use(['table'], function () {
    table = layui.table;
    layAjaxTable(table, {
        url: "/company/list",
        where: {flag: -1},
        cols: [[
            {title: '序号', width: 70, templet: '#indexTpl'}
            , {field: 'com_name', title: '名称', width: 200}
            , {templet: '#flagTpl', title: '启用状态', width: 100,align:"center"}
            , {field: 'com_desc', title: '备注', width: 100}
            , {title: '操作', toolbar: '#barDemo', width: 140}
        ]]
    });
    table.on('tool(demo)', function (obj) {
        var data = obj.data;
        var id = data["id"];
        if (obj.event === 'del') {
            layer.confirm('确认删除该数据？', {title: "系统提示"}, function (index) {
                ajaxDel(id, obj, index);
            });
        } else if (obj.event === 'edit') {
            juicer.renderTo("updatetpl", data, "update-form");
            form.render(null, "updateForm");
            editLayer();
        }
    });

});

function ajaxUpdate(data) {
    var url = "/company/update";
    data["flag"]=data["flag"]=="on"?1:0;
    postAjax(url, data, function (res) {
        if (res["code"] == 0) {
            layer.msg(res["msg"]);
            layer.close(dialogOpen);
            query();
        } else {
            layer.msg(res["msg"]);
        }

    }, function () {
        layer.msg("修改发生异常");
    });
}


function ajaxAdd(data) {
    var url = "/company/add";
    data["flag"]=data["flag"]=="on"?1:0;
    postAjax(url, data, function (res) {
        if (res["code"] == 0) {
            query();
            var msg = "<div style='width:100%;height: 50px;line-height: 50px; text-align: center;'>" + res["msg"] + "</div>";
            dialogGoon(msg, function () {

            }, function () {
                layer.close(dialogOpen);
            })
        } else {
            layer.alert(res["msg"]);
        }

    }, function () {
        layer.msg("添加发生异常");
    });
}

function ajaxDel(id, obj, index) {
    var url = "/company/del";
    postAjax(url, {"id": id}, function (res) {
        if (res["code"] == 0) {
            layer.msg(res["msg"]);
            obj.del();
            layer.close(index);
        } else {
            layer.msg(res["msg"]);
        }

    }, function () {
        layer.msg("删除发生异常");
    });
}

//添加
function listAddTo() {
    var dom = $("#add-form");
    $("#add-reset").click();
    dialogOpen = layer.open({
        title: '添加公司'
        , type: 1
        , area:['600px', "400px"]
        , content: dom
        , btn: ['确定', '取消']
        , yes: function () {
            $("#add-submit").click();
        }
    });
}

function editLayer() {
    var dom = $("#update-form");
    dialogOpen = layer.open({
        title: '编辑公司'
        , type: 1
        , area:['600px', "400px"]
        , content: dom
        , btn: ['确定', '取消']
        , yes: function () {
            $("#update-submit").click();
        }
    });
}

function query() {
    var name = $.trim($("#name-input").val());
    var state = $("#stateselect option:selected").val();
    table.reload('list-table', {where: {"com_name": name,"flag":state}});
}

$(function () {
    $("#add-btn").click(listAddTo);
    $("#query-btn").click(query);
    regEnterQuery(query)
});