var adminPage = {
    baseUrl: '/admins/',
    datagrid: $('#dg').datagrid({
        url: '/admins/',
        pagination: true,
        fitColumns: true,
        rownumbers: true,
        fit: true,
        toolbar: '#toolbar',
        method: 'get',
        singleSelect: true,
        columns: [[
            //{field: 'id', title: 'Id', sortable: true, width: 100,hidden:true},
            {field: 'admin_name', title: '用户名', sortable: true, width: 100},
            {
                field: 'create_time', title: '创建时间', sortable: true, width: 100,
                formatter: function (value, row, index) {
                    return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                }
            },
            {
                field: 'update_time', title: '最后一次更新时间', sortable: true, width: 100,
                formatter: function (value, row, index) {
                    return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                }
            }, {
                field: 'user',
                title: '是否有效',
                width: 100,
                formatter: function (value, row, index) {
                    if (value.invalid == 0) {
                        return "<img src='/static/jquery-easyui-1.5.1/themes/icons/cancel.png'>"
                    } else if (value.invalid == 1) {
                        return "<img src='/static/jquery-easyui-1.5.1/themes/icons/ok.png'>"
                    }
                }
            }, {
                field: 'id',
                title: '操作',
                width: 100,
                //hidden:true,
                formatter: function (val, row, index) {
                    var invalid = row.user.invalid;
                    var html = '';
                    if (invalid == 0) {
                        html += '<a href="#" onclick="adminPage.active(' + val + ',1)">激活</a>';
                    } else {
                        html += '<a href="#" onclick="adminPage.active(' + val + ',0)">注销</a>';
                    }
                    html += '|<a href="#" onclick="adminPage.grantRole(' + val + ')">授予角色</a>'
                    return html;

                }
            }
        ]]
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
                adminPage.form.submit();
            }
        }, {
            iconCls: 'icon-cancel',
            text: '取消',
            handler: function () {
                adminPage.dialog.dialog('close');
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
            if (r.ok) {
                adminPage.dialog.dialog('close');        // close the dialog
                adminPage.datagrid.datagrid('reload');    // reload the user data
            } else {
                $.showMsg(r.msg);
            }
            $.messager.progress('close');
        }
    }),

    queryDatagrid: function (queryParams) {
        this.datagrid.datagrid({
            queryParams: queryParams
        });
        return this;
    },

    add: function () {
        this.dialog.dialog('open').dialog('setTitle', '添加角色');
        this.form.form('clear');
        this.form.form({url: this.baseUrl});
        return this;
    },

    edit: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (row) {
            this.dialog.dialog('open').dialog('setTitle', '编辑角色');
            this.form.form('load', row);
            this.form.form({url: this.baseUrl + row.id});
        } else {
            $.alertMsg('请选择一条数据进行编辑！！！')
        }
        return this;
    },

    destroy: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (row) {
            var p = this;
            $.messager.confirm('系统提醒', '您确定删除这条数据吗?', function (r) {
                if (r) {
                    $.messager.progress();
                    $.ajax({
                        url: p.baseUrl + row.id,
                        type: 'DELETE',
                        dataType: 'json',
                        success: function (result) {
                            if (result.ok) {
                                p.datagrid.datagrid('reload');    // reload the user data
                            } else {
                                $.showMsg(result.msg)
                            }
                            $.messager.progress('close');
                        }
                    });
                }
            });
        } else {
            $.alertMsg('请选择一条数据进行删除！！！')
        }
    },

    active: function (id, invalid) {
        var msg = '您确定注销这条数据吗？';
        if (invalid == 1) {
            msg = '您确定激活这条数据吗？';
        }

        $.messager.confirm('系统提醒', msg, function (r) {
            if (r) {
                $.messager.progress();
                $.ajax({
                    url: '/adminValid/' + id + "/" + invalid,
                    type: 'PUT',
                    dataType: 'json',
                    success: function (result) {
                        if (result.ok) {
                            adminPage.datagrid.datagrid('reload');    // reload the user data
                        } else {
                            $.showMsg(result.msg)
                        }
                        $.messager.progress('close');
                    }
                });
            }
        });
    },

    grantRole: function(id) {
        var datagrid2 = $('#dg2').datagrid({
            url: '/roles//',
            fitColumns: true,
            rownumbers: true,
            fit: true,
            method: 'get',

            columns: [[
                //{field: 'id', title: 'Id', sortable: true, width: 100,hidden:true},
		        {field:'id', checkbox:true},
                {field: 'role_name', title: '角色名', sortable: true, width: 100},
                {
                    field: 'create_time', title: '创建时间', sortable: true, width: 100,
                    formatter: function (value, row, index) {
                        return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                    }
                }
	        ]]
        });
        var dlg2 = $('#dlg2').dialog({
            title:'授予角色',
            iconCls:'icon-save',
            buttons: [{
                text:'授权',
                iconCls:'icon-ok',
                handler:function(){
                    var roles = datagrid2.datagrid('getSelections')
                    console.info(roles)

                    var roleIds = new Array();
                    for(var i = 0; i < roles.length; i ++) {
                        roleIds[i] = roles[i].id;
                    }
                    $.ajax({
                        url: '/arrefs/',
                        type: 'POST',
                        dataType: 'json',
                        data : {adminId:id, roleIds:roleIds},
                        success: function (result) {
                            if (result.ok) {
                                dlg2.dialog('close');
                            }
                            $.showMsg(result.msg)
                        }
                    });
                }
            },{
                text:'取消',
                iconCls:'icon-cancel',
                handler:function(){
                    dlg2.dialog('close');
                }
            }]
        });
    }
};

