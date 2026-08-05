import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import './envfile.dart';

class EnvironmentEditor extends StatefulWidget {
  static const String zero = "";

  final String pluginId;
  final String raw;
  final void Function(String content)? onChange;

  const EnvironmentEditor(this.pluginId, this.raw, {super.key, this.onChange});

  static FutureBuilder<String> future(
    String pluginId,
    Future<String> pending, {
    void Function(String content)? onChange,
  }) {
    return ds.future(zero, pending, (snapshot) {
      return ds.ErrorScreen(
        EnvironmentEditor(
          pluginId,
          snapshot.data ?? zero,
          key: ValueKey(snapshot.data.hashCode),
          onChange: onChange,
        ),
        cause: snapshot.hasError ? ds.Error.unknown(snapshot.error!) : ds.Error.zero,
      );
    });
  }

  @override
  State<EnvironmentEditor> createState() => _EnvironmentEditorState();
}

class _EnvironmentEditorState extends State<EnvironmentEditor> {
  List<Variable> _variables = [];

  // tracked separately from _variables. _pendingGeneration only bumps on
  // submit (not per keystroke, unlike Variable.id/digest which change on
  // every edit) so the row's TextFormFields keep their focus while typing
  // and only remount - visibly clearing - once actually submitted.
  Variable _pending = Variable('', '', '');
  int _generation = 0;
  bool _dirty = false;

  @override
  void initState() {
    super.initState();
    _variables = parseEnv(widget.raw);
  }

  @override
  void dispose() {
    if (_dirty) _persist();
    super.dispose();
  }

  void _persist() {
    widget.onChange?.call(serializeEnv(_variables));
    _dirty = false;
  }

  void _persistIfDirty() {
    if (_dirty) _persist();
  }

  void _changePending({String? key, String? value}) {
    setState(() {
      _pending = _pending.copyWith(key: key, value: value);
    });
  }

  void _submitPending() {
    if (_pending.key.isEmpty) return;

    setState(() {
      _variables = [Variable(_pending.key, _pending.value, _pending.hint), ..._variables];
      _pending = Variable('', '', '');
      _generation++;
    });
    _persist();
  }

  void _change(Variable original, {String? key, String? value}) {
    setState(() {
      final idx = _variables.indexWhere((v) => identical(v, original));
      if (idx == -1) return;

      _variables[idx] = original.copyWith(key: key, value: value);
    });
    _dirty = true;
  }

  void _remove(Variable v) {
    setState(() {
      _variables.removeWhere((e) => identical(e, v));
    });
    _persist();
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return forms.Container(
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          forms.Field(
            key: ValueKey(_generation),
            input: Row(
              children: [
                Expanded(
                  child: TextFormField(
                    initialValue: _pending.key,
                    maxLines: 1,
                    decoration: const InputDecoration(hintText: "key"),
                    onChanged: (v) => _changePending(key: v),
                    onFieldSubmitted: (_) => _submitPending(),
                  ),
                ),
                Expanded(
                  child: TextFormField(
                    initialValue: _pending.value,
                    maxLines: 1,
                    decoration: const InputDecoration(hintText: "value"),
                    onChanged: (v) => _changePending(value: v),
                    onFieldSubmitted: (_) => _submitPending(),
                  ),
                ),
              ],
            ),
            trailing: [
              IconButton(
                icon: const Icon(Icons.add),
                tooltip: "add variable",
                onPressed: _submitPending,
              ),
            ],
          ),
          ..._variables.map(_field),
        ],
      ),
    );
  }

  Widget _field(Variable v) {
    return forms.Field(
      key: ValueKey(v.id),
      input: SelectionContainer.disabled(
        child: Row(
          children: [
            Expanded(
              child: Focus(
                onFocusChange: (hasFocus) {
                  if (!hasFocus) _persistIfDirty();
                },
                child: TextFormField(
                  initialValue: v.key,
                  maxLines: 1,
                  decoration: const InputDecoration(hintText: "key"),
                  onChanged: (updated) => _change(v, key: updated),
                ),
              ),
            ),
            Expanded(
              child: Focus(
                onFocusChange: (hasFocus) {
                  if (!hasFocus) _persistIfDirty();
                },
                child: TextFormField(
                  initialValue: v.value,
                  maxLines: 1,
                  decoration: const InputDecoration(hintText: "value"),
                  onChanged: (updated) => _change(v, value: updated),
                ),
              ),
            ),
          ],
        ),
      ),
      help: v.hint.isEmpty ? ds.HelpScope.None : ds.Hint(Text(v.hint)),
      trailing: [
        IconButton(
          icon: const Icon(Icons.close),
          tooltip: "remove",
          onPressed: () => _remove(v),
        ),
      ],
    );
  }
}
