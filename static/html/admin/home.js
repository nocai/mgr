function AdminPage() {
    var me = this;
    this.datagrid = $('#dg').datagrid({
        url: this.baseUrl,
        pagination: true,
        fitColumns: true,
        rownumbers: true,
        fit: true,
        toolbar: '#toolbar',
        method: 'get',
        singleSelect: true,
        columns: [[
            {field: 'id', title: 'Id', sortable: true},
            {field: 'admin_name', title: '用户名', width: 100, sortable: true},
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
            }, {
                field: 'invalid',
                title: '是否有效',
                width: 100,
                sortable: true,
                formatter: function (value, row, index) {
                    if (value == 1) {
                        return "<img src='/static/jquery-easyui-1.5.1/themes/icons/cancel.png'>"
                    } else if (value == 2) {
                        return "<img src='/static/jquery-easyui-1.5.1/themes/icons/ok.png'>"
                    }
                }
            }
        ]]
    });

    // 弹出窗口：添加or修改
    this.dialog = $('#dlg').dialog({
        width: 400,
        modal: true,
        closed: true,
        buttons:[{
            iconCls: 'icon-ok',
            text: '保存',
            handler: function () {
                me.form.submit();
            }
        }, {
            iconCls: 'icon-cancel',
            text: '取消',
            handler: function () {
                me.dialog.dialog('close');
            }
        }]
    }).dialog('center');

    // 表单：添加or修改
    this.form = $('#fm').form({
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
                me.dialog.dialog('close');        // close the dialog
                me.datagrid.datagrid('reload');    // reload the user data
            } else {
                $.showMsg(r.msg);
            }
            $.messager.progress('close');
        }
    });
}

AdminPage.prototype = {
    constructor: AdminPage,
    baseUrl: '/admins/',
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

    active: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (row) {
            if (row.invalid == 2) {
                $.showMsg('此帐号有效，无需激活...')
                return;
            }
            var p = this;
            $.messager.confirm('系统提醒', '您确定激活这条数据吗?', function (r) {
                if (r) {
                    $.messager.progress();
                    $.ajax({
                        url: '/adminValid/' + row.id + '/2',
                        type: 'PUT',
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
            $.alertMsg('请选择一条数据进行激活！！！')
        }
    }
};
