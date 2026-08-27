import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import './api.dart' as api;

class KnownMediaEdit extends StatelessWidget {
  final api.Known current;
  final forms.FnOnChange<api.Known> onChange;
  final Widget closable;
  final Widget deletable;
  final EdgeInsets? padding;

  const KnownMediaEdit(
    this.current, {
    super.key,
    this.onChange = forms.FnOnChangeNoop,
    this.closable = ds.Empty,
    this.deletable = ds.Empty,
    this.padding,
  });

  @override
  Widget build(BuildContext context) {
    return forms.Container(
      padding: padding,
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          forms.Field(label: Text("uid"), input: Text(current.uid)),
          forms.Field(
            label: Text("source id"),
            input: Text(current.id),
            trailing: [deletable, closable],
          ),
          forms.Field(
            label: Text("description"),
            input: TextFormField(
              readOnly: onChange == forms.FnOnChangeNoop,
              initialValue: current.description,
              onChanged: (v) => onChange(current..description = v),
            ),
          ),
          forms.Field(
            label: Text("summary"),
            input: TextFormField(
              readOnly: onChange == forms.FnOnChangeNoop,
              initialValue: current.summary,
              maxLines: null,
              onChanged: (v) => onChange(current..summary = v),
            ),
          ),
          forms.Field(
            label: Text("image"),
            input: TextFormField(
              readOnly: onChange == forms.FnOnChangeNoop,
              initialValue: current.image,
              onChanged: (v) => onChange(current..image = v),
            ),
          ),
          forms.Field(
            label: Text("released"),
            input: TextFormField(
              readOnly: onChange == forms.FnOnChangeNoop,
              initialValue: current.released,
              onChanged: (v) => onChange(current..released = v),
            ),
          ),
          forms.Field(
            label: Text("rating"),
            input: TextFormField(
              readOnly: onChange == forms.FnOnChangeNoop,
              initialValue: current.rating.toString(),
              keyboardType: TextInputType.numberWithOptions(decimal: true),
              onChanged: (v) => onChange(current..rating = double.tryParse(v) ?? current.rating),
            ),
          ),
          forms.Field(
            label: Text("adult"),
            input: forms.Checkbox(
              Text("adult content"),
              value: current.adult,
              onChanged: onChange == forms.FnOnChangeNoop ? null : (v) => onChange(current..adult = v ?? false),
            ),
          ),
          forms.Field(label: Text("mimetype"), input: Text(current.mimetype)),
          forms.Field(label: Text("source"), input: Text(current.source)),
        ],
      ),
    );
  }
}
