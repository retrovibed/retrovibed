import 'package:flutter/material.dart';
import 'package:pdfrx/pdfrx.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import './api.dart' as api;
import './media.pb.dart';
import './preview.unsupported.dart';

// pdfrx fetches the document over the same authenticated /m/{id} route the rest of the
// library reads through, and requests the byte ranges it needs per page rather than the
// whole file.
class PreviewPdf extends StatelessWidget {
  final Media current;

  const PreviewPdf({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    final uri = api.media.download_uri(current.id);

    return PdfViewer.uri(
      Uri.parse(uri),
      headers: httpx.localheaders(uri) ?? const {},
      params: PdfViewerParams(
        errorBannerBuilder: (context, error, stack, documentRef) => PreviewUnsupported(
          current: current,
          description: const Text("unable to render this document"),
        ),
        loadingBannerBuilder: (context, bytes, total) => ds.Loading(ds.Empty, loading: true),
      ),
    );
  }
}
