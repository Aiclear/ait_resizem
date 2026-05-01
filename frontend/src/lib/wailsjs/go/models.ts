export namespace rimage {
	
	export class RenameOptions {
	    enabled: boolean;
	    template: string;
	    start_index: number;
	    index_padding: number;
	    date_format: string;
	
	    static createFrom(source: any = {}) {
	        return new RenameOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.template = source["template"];
	        this.start_index = source["start_index"];
	        this.index_padding = source["index_padding"];
	        this.date_format = source["date_format"];
	    }
	}
	export class WatermarkOptions {
	    type: number;
	    text: string;
	    font_path: string;
	    font_size: number;
	    font_color: string;
	    image_path: string;
	    opacity: number;
	    position: number;
	    offset_x: number;
	    offset_y: number;
	    rotation: number;
	    scale: number;
	
	    static createFrom(source: any = {}) {
	        return new WatermarkOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.text = source["text"];
	        this.font_path = source["font_path"];
	        this.font_size = source["font_size"];
	        this.font_color = source["font_color"];
	        this.image_path = source["image_path"];
	        this.opacity = source["opacity"];
	        this.position = source["position"];
	        this.offset_x = source["offset_x"];
	        this.offset_y = source["offset_y"];
	        this.rotation = source["rotation"];
	        this.scale = source["scale"];
	    }
	}
	export class ImageOptions {
	    dest_format: number;
	    resample_filter: number;
	    desc_path: string;
	    dest_width: number;
	    dest_height: number;
	    jpeg_quality: number;
	    gif_number_of_colors: number;
	    tiff_compression: number;
	    png_compression: number;
	    auto_orientation: boolean;
	    cpu_memory_usage: number;
	    watermark: WatermarkOptions;
	    rename: RenameOptions;
	
	    static createFrom(source: any = {}) {
	        return new ImageOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dest_format = source["dest_format"];
	        this.resample_filter = source["resample_filter"];
	        this.desc_path = source["desc_path"];
	        this.dest_width = source["dest_width"];
	        this.dest_height = source["dest_height"];
	        this.jpeg_quality = source["jpeg_quality"];
	        this.gif_number_of_colors = source["gif_number_of_colors"];
	        this.tiff_compression = source["tiff_compression"];
	        this.png_compression = source["png_compression"];
	        this.auto_orientation = source["auto_orientation"];
	        this.cpu_memory_usage = source["cpu_memory_usage"];
	        this.watermark = this.convertValues(source["watermark"], WatermarkOptions);
	        this.rename = this.convertValues(source["rename"], RenameOptions);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class SupportedOutputImageType {
	    name: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new SupportedOutputImageType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class SupportedResampleFilterType {
	    name: string;
	    value: number;
	
	    static createFrom(source: any = {}) {
	        return new SupportedResampleFilterType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}

}

