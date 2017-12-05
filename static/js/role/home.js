var rolePage = {
    baseUrl: '/roles/',
    datagrid: $('#dg').datagrid({
        //url: this.baseUrl,
        url: '/roles/',
        pagination: true,
        fitColumns: true,
        rownumbers: true,
        fit: true,
        toolbar: '#toolbar',
        method: 'get',
        singleSelect: true,
        columns: [[
            {field: 'id', title: 'Id', sortable: true},
            {field: 'role_name', title: '角色名', width: 100, sortable: true},
            {
                field: 'create_time', title: '创建时间', width: 100, sortable: true,
                formatter: function (value, row, index) {
                    return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                }
            },
            {
                field: 'update_time', title: '最后一次更新时间', width: 100, sortable: true,
                formatter: function (value, row, index) {
                    return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                }
            }
        ]],
        // loader:function(_253,_254,_255){
        //     alert('a')
        //     var opts=$(this).datagrid("options");
        //     if(!opts.url){
        //         return false;
        //     }
        //     $.ajax({type:opts.method,url:opts.url,data:_253,dataType:"json",success:function(data){
        //         console.info(data)
        //
        //             _254(data.data);
        //         },error:function(){
        //             _255.apply(this,arguments);
        //         }}
        //     );
        // }

    }),

    // 弹出窗口：添加or修改
    dialog: $('#dlg').dialog({
        width: 400,
        modal: true,
        closed: true,
        buttons: [{
            iconCls: 'icon-ok',
            text: '保存',
            handler: function () {
                rolePage.form.submit();
            }
        }, {
            iconCls: 'icon-cancel',
            text: '取消',
            handler: function () {
                rolePage.dialog.dialog('close');
            }
        }]
    }).dialog('center'),

    // 表单：添加or修改
    form: $('#fm').form({
        onSubmit: function () {
            var valid = $(this).form('validate');
            if (valid) {
                $.messager.progress();
                return true;
            }
            return false;
        },
        success: function (r) {
            r = eval('(' + r + ')');
            if (r.code === 0) {
                rolePage.dialog.dialog('close');        // close the dialog
                rolePage.datagrid.datagrid('reload');    // reload the user data
            } else {
                $.showMsg(r.msg);
            }
            $.messager.progress('close');
        }
    }),

    add: function () {
        this.dialog.dialog('open').dialog('setTitle', '添加角色');
        this.form.form('clear');
        this.form.form({url: this.baseUrl});
        return this;
    },

    edit: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (!row) {
            $.alertMsg('请选择一条数据进行编辑！！！')
            return this;
        }
        this.dialog.dialog('open').dialog('setTitle', '编辑角色');
        this.form.form('load', row);
        this.form.form({url: this.baseUrl + row.id});
        return this;
    },

    destroy: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (!row) {
            $.alertMsg('请选择一条数据进行删除！！！');
            return this;
        }
        $.messager.confirm('系统提醒', '您确定删除这条数据吗?', function (r) {
            if (r) {
                $.messager.progress();
                $.ajax({
                    url: rolePage.baseUrl + row.id,
                    type: 'DELETE',
                    dataType: 'json',
                    success: function (result) {
                        if (result.code === 0) {
                            rolePage.datagrid.datagrid('reload');    // reload the user data
                        } else {
                            $.showMsg(result.msg)
                        }
                        $.messager.progress('close');
                    }
                });
            }
        });
    }

};

