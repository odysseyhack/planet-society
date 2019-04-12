//
//  TransactionTableViewCell.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

struct TransactionTableViewCellViewModel {
    let image: UIImage
    let title: String
    let subtitle: String
}

final class TransactionTableViewCell: UITableViewCell {

    // MARK: - Static properties

    private let stackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.spacing = 20
        stackView.alignment = .center

        stackView.backgroundColor = .red

        return stackView
    }()

    private let itemImageView: UIImageView = {

        let imageView = UIImageView()
        imageView.translatesAutoresizingMaskIntoConstraints = false

        let dimension: CGFloat = 30
        imageView.widthAnchor.constraint(equalToConstant: dimension).isActive = true
        imageView.heightAnchor.constraint(equalToConstant: dimension).isActive = true

        return imageView
    }()

    private let itemTitleLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.regular(ofSize: 12)
        label.textColor = PHColors.greyishBrown

        return label
    }()

    private let itemSubtitleLabel: UILabel = {

        let label = UILabel()
        label.font = PHFonts.regular(ofSize: 11)
        label.textColor = PHColors.grey

        return label
    }()

    private let verticalStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .vertical
        stackView.spacing = 2
        stackView.alignment = .leading

        return stackView
    }()

    private let infoButton: UIButton = {

        let button = UIButton()
        button.translatesAutoresizingMaskIntoConstraints = false

        let image = UIImage(named: "info_button")
        button.setImage(image, for: .normal)

        return button
    }()

    private let checkmarkImageView: UIImageView = {

        let imageView = UIImageView()
        imageView.translatesAutoresizingMaskIntoConstraints = false

        let image = UIImage(named: "checkmark")
        imageView.image = image

        return imageView
    }()

    private lazy var separatorView: UIView = {

        let view = UIView()
        view.translatesAutoresizingMaskIntoConstraints = false

        view.backgroundColor = PHColors.grey.withAlphaComponent(0.5)

        return view
    }()

    // MARK: - Initialization

    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: style, reuseIdentifier: reuseIdentifier)

        configure()
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Configuration

    private func configure() {

        addSubview(stackView)
        addSubview(separatorView)

        let margin: CGFloat = 15
        stackView.topAnchor.constraint(equalTo: topAnchor, constant: margin).isActive = true
        stackView.leftAnchor.constraint(equalTo: leftAnchor, constant: margin).isActive = true
        stackView.rightAnchor.constraint(equalTo: rightAnchor, constant: -margin).isActive = true

        separatorView.topAnchor.constraint(equalTo: stackView.bottomAnchor, constant: margin).isActive = true
        separatorView.leftAnchor.constraint(equalTo: leftAnchor).isActive = true
        separatorView.rightAnchor.constraint(equalTo: rightAnchor).isActive = true
        separatorView.bottomAnchor.constraint(equalTo: bottomAnchor).isActive = true
        separatorView.heightAnchor.constraint(equalToConstant: 1).isActive = true
        
        stackView.addArrangedSubview(itemImageView)
        stackView.addArrangedSubview(verticalStackView)
        verticalStackView.addArrangedSubview(itemTitleLabel)
        verticalStackView.addArrangedSubview(itemSubtitleLabel)
        stackView.addArrangedSubview(infoButton)
        stackView.addArrangedSubview(checkmarkImageView)
    }

    func configure(withViewModel viewModel: TransactionTableViewCellViewModel) {

        itemTitleLabel.text = viewModel.title
        itemSubtitleLabel.text = viewModel.subtitle
    }
}
