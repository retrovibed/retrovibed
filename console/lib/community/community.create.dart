import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'community.pb.dart';
import 'community.edit.dart';

class CommunityCreate extends StatefulWidget {
  final Future<CommunityCreateResponse> Function(Community) create;
  final void Function(Community) onCreate;
  final VoidCallback onCancel;
  final BoxConstraints? constraints;
  final EdgeInsets? margin;
  final EdgeInsets? padding;

  const CommunityCreate({
    super.key,
    required this.create,
    required this.onCreate,
    required this.onCancel,
    this.constraints,
    this.margin,
    this.padding,
  });

  @override
  _CommunityCreateState createState() => _CommunityCreateState();
}

class _CommunityCreateState extends State<CommunityCreate> with ds.LoadingState {
  Community _community = Community(mimetype: mimex.binary);
  bool _creating = false;

  void _createCommunity() {
    setState(() {
      _creating = true;
      cause = ds.Error.zero;
    });

    widget
        .create(_community)
        .then((response) => widget.onCreate(response.community))
        .catchError((c) {
          setState(() {
            cause = ds.Error.conflict(
              c,
              onTap: reseterr,
              message: Text("a community with this url already exists"),
            );
          });
        }, test: httpx.ErrorsTest.conflict)
        .catchError((c) {
          setState(() {
            cause = ds.Errors.httpauto(c, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.httpauto)
        .catchError((c) {
          setState(() {
            cause = ds.Error.unknown(c, onTap: reseterr);
          });
        })
        .whenComplete(() {
          setState(() {
            _creating = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Container(
      padding: widget.padding ?? defaults.padding,
      margin: widget.margin ?? defaults.margin,
      constraints: widget.constraints,
      Column(
        spacing: defaults.spacing,
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            'New Community',
            style: theme.textTheme.headlineSmall,
            textAlign: TextAlign.center,
          ),
          ds.Loading(
            cause: cause,
            CommunityEdit(
              community: _community,
              onChange: (c) => setState(() => _community = c),
              autofocus: true,
            ),
          ),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextButton(
                onPressed: _creating ? null : widget.onCancel,
                child: Text('Cancel'),
              ),
              ElevatedButton(
                onPressed: _creating ? null : _createCommunity,
                child: _creating
                    ? SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : Text('Create'),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
