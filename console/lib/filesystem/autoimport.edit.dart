import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'autoimport.api.dart' as api;

// edits a monitored directory in place, invoking onChange with the updated record. the daemon
// treats the path as immutable once created so it is only editable while creating.
class AutoImportEdit extends StatelessWidget {
  final api.AutoimportDirectory current;
  final Function(api.AutoimportDirectory) onChange;
  final bool pathEditable;
  // buttons rendered below the fields, e.g. save and delete.
  final List<Widget> actions;

  AutoImportEdit({
    super.key,
    api.AutoimportDirectory? current,
    this.onChange = ds.fnNoop,
    this.pathEditable = false,
    this.actions = const [],
  }) : current = current ?? api.autoimport.directory();

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          forms.Field(
            label: const Text("path"),
            input: TextFormField(
              initialValue: current.path,
              maxLines: 1,
              readOnly: !pathEditable,
              enabled: pathEditable,
              autofocus: pathEditable,
              onChanged: (v) => onChange(current..path = v),
            ),
            trailing: actions,
          ),
          forms.Field(
            label: const Text("description"),
            input: TextFormField(
              initialValue: current.description,
              maxLines: 1,
              onChanged: (v) => onChange(current..description = v),
            ),
          ),
          forms.Field(
            label: const Text("debounce (minutes)"),
            input: TextFormField(
              decoration: const InputDecoration(
                helperText: "files are imported once they haven't been modified for this long",
              ),
              initialValue: (current.debounce.toInt() ~/ 60).toString(),
              keyboardType: TextInputType.number,
              onChanged: (v) {
                final minutes = int.tryParse(v);
                if (minutes == null) return;
                onChange(current..debounce = ds.Int64(minutes * 60));
              },
            ),
          ),
          forms.Field(
            label: const Text("mode"),
            input: DropdownButtonFormField<int>(
              initialValue: current.mode,
              items: [
                DropdownMenuItem(
                  value: api.autoimportModeCopy,
                  child: Text(api.autoimportModeLabel(api.autoimportModeCopy)),
                ),
                DropdownMenuItem(
                  value: api.autoimportModeMove,
                  child: Text(api.autoimportModeLabel(api.autoimportModeMove)),
                ),
              ],
              onChanged: (v) => onChange(current..mode = v ?? current.mode),
            ),
          ),
          forms.Field(
            label: const Text("library directory"),
            input: TextFormField(
              decoration: const InputDecoration(helperText: "blank imports into the library root"),
              initialValue: current.libraryDirectoryId,
              maxLines: 1,
              onChanged: (v) => onChange(current..libraryDirectoryId = v),
            ),
          ),
        ],
      ),
    );
  }
}
