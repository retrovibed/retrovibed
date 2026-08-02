import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'community.pb.dart';
import 'community.edit.dart';

class CommunityUpdate extends StatefulWidget {
  final Community community;
  final Future<CommunityUpdateResponse> Function(Community) update;
  final void Function(Community) onUpdate;
  final VoidCallback onCancel;
  final EdgeInsets? padding;
  final EdgeInsets? margin;
  final BoxDecoration? decoration;
  final BoxConstraints? constraints;
  final Alignment? alignment;
  final Clip clipBehavior;

  const CommunityUpdate({
    super.key,
    required this.community,
    required this.update,
    required this.onUpdate,
    required this.onCancel,
    this.padding,
    this.margin,
    this.decoration,
    this.constraints,
    this.alignment,
    this.clipBehavior = Clip.none,
  });

  @override
  _CommunityUpdateState createState() => _CommunityUpdateState();
}

class _CommunityUpdateState extends State<CommunityUpdate> {
  Community _community = Community();
  Widget _cause = ds.Error.zero;

  @override
  void initState() {
    super.initState();
    _community = widget.community.deepCopy();
  }

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _clearCause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  Future<void> _save() {
    setState(() => _cause = ds.Error.zero);

    return widget
        .update(_community)
        .then((response) => widget.onUpdate(response.community))
        .catchError((cause) {
          setState(() {
            _cause = ds.Errors.httpauto(cause, onTap: _clearCause);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _clearCause);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      padding: widget.padding ?? defaults.padding,
      margin: widget.margin,
      decoration: widget.decoration ?? const BoxDecoration(),
      constraints: widget.constraints,
      alignment: widget.alignment,
      clipBehavior: widget.clipBehavior,
      Column(
        spacing: defaults.spacing,
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          ds.Loading(
            cause: _cause,
            CommunityEdit(
              community: _community,
              onChange: (c) => setState(() => _community = c),
              readOnly: true,
              autofocus: true,
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextButton(
                onPressed: widget.onCancel,
                child: Text('Cancel'),
              ),
              ds.LoadingButton(
                Text('Save'),
                onPressed: _save,
              ),
            ],
          ),
        ],
      ),
    );
  }
}
