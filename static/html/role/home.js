function RolePage() {
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
    });
    this.dialog = $('#dlg').dialog({
        width:400,
        modal: true,
        closed: true,
        buttons: '#dlg-buttons',
    }).dialog('center');

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
                page.dialog.dialog('close');        // close the dialog
                page.datagrid.datagrid('reload');    // reload the user data
            } else {
                $.showMsg(r.msg);
            }
            $.messager.progress('close');
        }
    });
}

RolePage.prototype = {
    constructor: RolePage,
    baseUrl: '/roles/',
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
        }
        return this;
    },

    save: function () {
        this.form.submit();
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
                                $.showMsg(result.msg)
                            }
                            $.messager.progress('close');
                        }
                    });
                }
            });
        }
    }
};
