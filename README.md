<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/mkubasz/quanto-magis">
    <img src="./docs/public/logo.jpeg" alt="Logo" width="160" height="160">
  </a>

<h3 align="center">Quanto Magis</h3>

  <p align="center">
    Spark engine for data processing based on Go using Kubernetes
    <br />
    <a href="https://github.com/mkubasz/quanto-magis"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/mkubasz/quanto-magis">View Demo</a>
    ·
    <a href="https://github.com/mkubasz/quanto-magis/issues">Report Bug</a>
    ·
    <a href="https://github.com/mkubasz/quanto-magis/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

This project is a simple example of how to use Go to create a Spark-like engine for data processing. The main goal is to create a simple and fast engine for data processing using Kubernetes. The project is in the early stage of development and is not ready for production use.



<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* [![Go Lang][Go]][Go-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* GO website => https://go.dev/doc/install

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/mkubasz/quanto-magis.git
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

```go
	session := NewQuantoSession().GetOrCreate()
	data := []interface{}{[]interface{}{"A", "B", "C"}, []interface{}{"D", "E"}}

	rdd := session.Parallelize(data)
	result := rdd.FlatMap(lowerCase)
    // result = [a, b, c, d, e]
```


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- Features -->
## Features

- Quanto Session - a session for data processing for a specific core
- RDD - Resilient Distributed Dataset
- DataFrame - a distributed collection of data organized into named columns

See the [open issues](https://github.com/mkubasz/quanto-magis/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- Architecture -->
## Architecture

[Link to architecture](https://link.excalidraw.com/l/6OWSHwUWoad/2foryRzJ0RH)

![Architecture](./docs/Quanto%20Magis%20lib.png "Title")

<!-- NOTES -->
## Notes
My journey learning tech.

[GO Notes](docs/NOTES.md)

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Mateusz Kubaszek - [@MateuszKubaszek](https://twitter.com/MateuszKubaszek) - mkubasz@gmail.com

Project Link: [https://github.com/mkubasz/quanto-magis](https://github.com/mkubasz/quanto-magis)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/mkubasz/quanto-magis.svg?style=for-the-badge
[contributors-url]: https://github.com/mkubasz/quanto-magis/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/mkubasz/quanto-magis.svg?style=for-the-badge
[forks-url]: https://github.com/mkubasz/quanto-magis/network/members
[stars-shield]: https://img.shields.io/github/stars/mkubasz/quanto-magis.svg?style=for-the-badge
[stars-url]: https://github.com/mkubasz/quanto-magis/stargazers
[issues-shield]: https://img.shields.io/github/issues/mkubasz/quanto-magis.svg?style=for-the-badge
[issues-url]: https://github.com/mkubasz/quanto-magis/issues
[license-shield]: https://img.shields.io/github/license/mkubasz/quanto-magis.svg?style=for-the-badge
[license-url]: https://github.com/mkubasz/quanto-magis/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/linkedin_username
[product-screenshot]: images/screenshot.png
[Go]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/

