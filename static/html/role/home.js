/**
 * Created by liujun on 17/3/27.
 */
function RolePage(datagrid_id, dialog_id, form_id) {
    this.datagrid_id = datagrid_id;
    this.datagrid = $('#' + this.datagrid_id).datagrid({
        url: this.baseUrl,
        pagination: true,
        fitColumns: true,
        rownumbers: true,
        fit: true,
        toolbar: '#toolbar',
        method: 'get',
        singleSelect: true
    });

    this.dialog_id = dialog_id;
    this.dialog = $('#' + this.dialog_id).dialog('center');
    this.form_id = form_id;
    this.form = $('#' + this.form_id).form({
        url: this.url
    });
}


RolePage.prototype = {
    constructor: RolePage,
    baseUrl: '/roles/',
    loadDatagrid: function () {
        this.datagrid.datagrid({
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
            ]]
        })
    },

    queryDatagrid: function (queryParams) {
        console.debug(queryParams);
        this.datagrid.datagrid({
            queryParams: queryParams
        });
        return this;
    },

    add: function () {
        this.dialog.dialog('open').dialog('setTitle', '添加角色');
        $('#fm').form('clear');
        this.url = this.baseUrl;
    },

    edit: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (row) {
            this.dialog.dialog('open').dialog('setTitle', '编辑角色');
            $('#fm').form('load', row);
            this.url = this.baseUrl + row.id;
        }
    },

    save: function () {
        var p = this;
        this.form.form('submit', {
            url: this.url,
            onSubmit: function () {
                if ($(this).form('validate')) {
                    $.messager.progress();
                    return true;
                }
                return false;
            },
            success: function (r) {
                r = eval('(' + r + ')');
                if (r.ok) {
                    p.dialog.dialog('close');        // close the dialog
                    p.datagrid.datagrid('reload');    // reload the user data
                } else {
                    $.messager.show({
                        title: '系统提示',
                        msg: r.msg
                    });
                }
                $.messager.progress('close');
            }
        });
    },

    destroy: function () {
        var row = $('#dg').datagrid('getSelected');
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
                                $.messager.show({    // show error message
                                    title: 'Error',
                                    msg: result.msg
                                });
                            }
                            $.messager.progress('close');
                        }
                    });
                }
            });
        }
    }
};
