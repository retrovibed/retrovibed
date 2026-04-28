import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/uuidx.dart' as uuidx;
import 'package:retrovibed/library/api.dart' as library_api;
import 'package:retrovibed/media/api.dart' as media_api;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/media/media.known.pb.dart';

class PublishMetadata extends StatefulWidget {
  final Download? download;
  final void Function(Known) onConfirm;
  final Future<KnownLookupResponse> Function(String, {List<httpx.Option> options}) knownGet;
  final Future<KnownCreateResponse> Function(KnownCreateRequest, {List<httpx.Option> options}) knownCreate;
  final Future<MetadataSyncResponse> Function(String, Media, {List<httpx.Option> options}) metadataSync;

  const PublishMetadata({
    super.key,
    required this.download,
    required this.onConfirm,
    this.knownGet = library_api.known.get,
    this.knownCreate = library_api.known.create,
    this.metadataSync = media_api.discovered.metadatasync,
  });

  @override
  State<PublishMetadata> createState() => _PublishMetadataState();
}

class _PublishMetadataState extends State<PublishMetadata> {
  Known _formData = Known(
    released: DateTime.now().toIso8601String().split('T').first,
    adult: false,
  );
  String _dirty = uuidx.min();
  bool _loadingMetadata = false;
  bool _hasExistingMetadata = false;
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _clearCause() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void initState() {
    super.initState();
    _initializeForm();
  }

  void _initializeForm() {
    final media = widget.download?.media ?? Media();
    final knownMediaId = media.knownMediaId;
    if (knownMediaId.isEmpty || uuidx.isMinMax(uuidx.fromString(knownMediaId))) {
      setState(() {
        _formData = Known(
          description: media.description,
          released: media.createdAt,
          adult: false,
        );
        _dirty = uuidx.random();
      });
      return;
    }

    setState(() {
      _loadingMetadata = true;
      _cause = ds.Error.zero;
    });

    final authOptions = [authn.AuthzCache.bearer(context)];
    widget
        .knownGet(knownMediaId, options: [httpx.Accept.json, ...authOptions])
        .then((response) {
          setState(() {
            _formData = response.known;
            _hasExistingMetadata = true;
            _dirty = uuidx.random();
          });
        })
        .catchError((cause) {
          setState(() {
            _formData = Known(
              description: media.description,
              released: media.createdAt,
              adult: false,
            );
            _dirty = uuidx.random();
          });
        })
        .whenComplete(() {
          setState(() {
            _loadingMetadata = false;
          });
        });
  }

  void _update(void Function(Known) fn) {
    fn(_formData);
    setState(() {});
  }

  bool _hasTorrent() {
    return !uuidx.isMinMax(
      uuidx.fromString(widget.download?.media.torrentId ?? uuidx.min()),
    );
  }

  void _sync(Known known) {
    if (!_hasTorrent()) {
      widget.onConfirm(known);
      return;
    }

    setState(() {
      _loadingMetadata = true;
      _cause = ds.Error.zero;
    });

    final authOptions = [authn.AuthzCache.bearer(context)];
    final updatedMedia = widget.download!.media..knownMediaId = known.id;

    httpx
        .withRetry(
          () => widget.metadataSync(
            widget.download!.media.torrentId,
            updatedMedia,
            options: authOptions,
          ),
        )
        .then((_) => widget.onConfirm(known))
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _clearCause);
          });
        })
        .whenComplete(() {
          setState(() {
            _loadingMetadata = false;
          });
        });
  }

  void _submit() {
    if (_formData.description.isEmpty) {
      setState(() {
        _cause = ds.Error.unknown('Title is required', onTap: _clearCause);
      });
      return;
    }

    if (_hasExistingMetadata) {
      _sync(_formData);
      return;
    }

    setState(() {
      _loadingMetadata = true;
      _cause = ds.Error.zero;
    });

    final authOptions = [authn.AuthzCache.bearer(context)];
    final req = KnownCreateRequest(known: _formData);

    widget
        .knownCreate(
          req,
          options: [httpx.Accept.json, httpx.Content.json, ...authOptions],
        )
        .then((response) => _sync(response.known))
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: _clearCause);
          });
        })
        .whenComplete(() {
          setState(() {
            _loadingMetadata = false;
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);
    final isLoading = _loadingMetadata;
    final content = forms.Container(
      padding: EdgeInsets.zero,
      decoration: BoxDecoration(borderRadius: defaults.borderRadius),
      Column(
        mainAxisSize: MainAxisSize.min,
        spacing: defaults.spacing,
        children: [
          Text(
            _hasExistingMetadata ? 'Review' : 'Metadata',
            style: theme.textTheme.bodyMedium,
            textAlign: TextAlign.center,
          ),
          Column(
            key: ValueKey(_dirty),
            mainAxisSize: MainAxisSize.min,
            children: [
              forms.Field(
                label: Text('Title'),
                input: TextFormField(
                  initialValue: _formData.description,
                  decoration: InputDecoration(
                    hintText: 'Content title',
                    border: OutlineInputBorder(),
                  ),
                  onChanged: (v) => _update((k) => k.description = v),
                ),
              ),
              forms.Field(
                label: Text('Summary'),
                input: TextFormField(
                  initialValue: _formData.summary,
                  decoration: InputDecoration(
                    hintText: 'Content summary',
                    border: OutlineInputBorder(),
                  ),
                  maxLines: 3,
                  onChanged: (v) => _update((k) => k.summary = v),
                ),
              ),
              forms.Field(
                label: Text('Image URL (optional)'),
                input: TextFormField(
                  initialValue: _formData.image,
                  decoration: InputDecoration(
                    hintText: 'https://example.com/image.jpg',
                    border: OutlineInputBorder(),
                  ),
                  onChanged: (v) => _update((k) => k.image = v),
                ),
              ),
              forms.Field(
                label: Text('Release Date'),
                input: TextFormField(
                  initialValue: _formData.released.split('T').first,
                  decoration: InputDecoration(
                    hintText: 'YYYY-MM-DD',
                    border: OutlineInputBorder(),
                  ),
                  onChanged: (v) => _update((k) => k.released = v),
                ),
              ),
              forms.Checkbox(
                Text('Adult content'),
                value: _formData.adult,
                onChanged: (v) => _update((k) => k.adult = v ?? false),
                description: Text('Mark if this content is for adults only'),
              ),
            ],
          ),
          Center(
            child: ElevatedButton.icon(
              onPressed: isLoading ? null : _submit,
              icon: Icon(Icons.check),
              label: Text('Continue'),
            ),
          ),
        ],
      ),
    );

    return ds.Loading(loading: isLoading, cause: _cause, content);
  }
}
