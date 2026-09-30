import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'autoimport.api.dart' as api;

// edits a monitored directory in place, invoking onChange with the updated record. the daemon
// treats the path as immutable once created so it is only editable while creating.
class AutoimportEdit extends StatelessWidget {
  final api.AutoimportDirectory current;
  final Function(api.AutoimportDirectory)? onChange;
  final bool pathEditable;

  AutoimportEdit({super.key, api.AutoimportDirectory? current, this.onChange, this.pathEditable = false})
    : current = current ?? api.autoimport.directory();

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          TextFormField(
            decoration: const InputDecoration(helperText: "path"),
            initialValue: current.path,
            readOnly: !pathEditable,
            enabled: pathEditable,
            autofocus: pathEditable,
            onChanged: (v) => onChange?.call(current..path = v),
          ),
          TextFormField(
            decoration: const InputDecoration(helperText: "description"),
            initialValue: current.description,
            onChanged: (v) => onChange?.call(current..description = v),
          ),
          TextFormField(
            decoration: const InputDecoration(helperText: "debounce (minutes)"),
            initialValue: (current.debounce.toInt() ~/ 60).toString(),
            keyboardType: TextInputType.number,
            onChanged: (v) {
              final minutes = int.tryParse(v);
              if (minutes == null) return;
              onChange?.call(current..debounce = ds.Int64(minutes * 60));
            },
          ),
          DropdownButtonFormField<int>(
            decoration: const InputDecoration(helperText: "mode"),
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
            onChanged: (v) => onChange?.call(current..mode = v ?? current.mode),
          ),
          TextFormField(
            decoration: const InputDecoration(helperText: "library directory (blank for the root)"),
            initialValue: current.libraryDirectoryId,
            onChanged: (v) => onChange?.call(current..libraryDirectoryId = v),
          ),
        ],
      ),
    );
  }
}
